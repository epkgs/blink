package downloader

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"net/http/cookiejar"
	netUrl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/epkgs/blink/internal/log"
	"github.com/jlaffaye/ftp"
	"github.com/lxn/win"
)

type IBeforeDownloadInterceptor func(job *Job)
type IHttpDownloadingInterceptor func(resp *http.Response, req *http.Request) io.Reader
type IFtpDownloadingInterceptor func(resp *ftp.Response, url netUrl.URL) io.Reader
type IBeforeSaveFileInterceptor func(job *Job)

type IInterceptors struct {
	BeforeDownload  IBeforeDownloadInterceptor
	HttpDownloading IHttpDownloadingInterceptor
	FtpDownloading  IFtpDownloadingInterceptor
	BeforeSaveFile  IBeforeSaveFileInterceptor
}

type Option struct {
	Dir                  string         // 下载路径，如果为空则使用当前目录
	FileNamePrefix       string         // 文件名前缀，默认空
	MaxThreads           int            // 下载线程，默认4
	MinChunkSize         uint64         // 最小分块大小，默认500KB
	EnableSaveFileDialog bool           // 是否打开保存文件对话框，默认false
	Overwrite            bool           // 是否覆盖已存在的文件，默认false
	Timeout              time.Duration  // 超时时间，默认10秒
	InsecureSkipVerify   bool           // 跳过证书验证，默认false
	Cookies              []*http.Cookie // 请求头Cookie，默认空。

	Interceptors IInterceptors
}

type Downloader struct {
	lastJobId uint64
	Option
}

type Job struct {
	downloader *Downloader
	Option

	id             uint64
	Url            *netUrl.URL
	FileName       string
	FileSize       uint64
	isSupportRange bool
	isFtp          bool
	httpClient     *http.Client

	_lck *sync.Mutex // 用于确保写入文件的顺序
}

func (opt Option) cloneOption() Option {
	return Option{
		Dir:                  opt.Dir,
		FileNamePrefix:       opt.FileNamePrefix,
		MaxThreads:           opt.MaxThreads,
		MinChunkSize:         opt.MinChunkSize,
		EnableSaveFileDialog: opt.EnableSaveFileDialog,
		Overwrite:            opt.Overwrite,
		Timeout:              opt.Timeout,
		InsecureSkipVerify:   opt.InsecureSkipVerify,
		Cookies:              opt.Cookies,

		Interceptors: IInterceptors{
			BeforeDownload:  opt.Interceptors.BeforeDownload,
			HttpDownloading: opt.Interceptors.HttpDownloading,
			FtpDownloading:  opt.Interceptors.FtpDownloading,
			BeforeSaveFile:  opt.Interceptors.BeforeSaveFile,
		},
	}
}

func New(withOption ...func(*Option)) *Downloader {

	// 默认参数
	opt := Option{
		Dir:                  "",
		FileNamePrefix:       "",
		MaxThreads:           4,
		MinChunkSize:         500 * 1024, // 500KB
		EnableSaveFileDialog: false,
		Overwrite:            false,
		Timeout:              10 * time.Second,
		InsecureSkipVerify:   false,
		Cookies:              make([]*http.Cookie, 0),

		Interceptors: IInterceptors{
			BeforeDownload: func(job *Job) {}, // 默认空实现
			HttpDownloading: func(resp *http.Response, req *http.Request) io.Reader {
				return resp.Body
			},
			FtpDownloading: func(resp *ftp.Response, url netUrl.URL) io.Reader {
				return resp
			},
			BeforeSaveFile: func(job *Job) {}, // 默认空实现
		},
	}

	for _, set := range withOption {
		set(&opt)
	}

	downloader := &Downloader{
		lastJobId: 0,
		Option:    opt,
	}

	return downloader
}

func (d *Downloader) Download(url string, withOption ...func(*Option)) (string, error) {
	job, err := d.NewJob(url, withOption...)
	if err != nil {
		return "", err
	}
	return job.Download()
}

func (d *Downloader) NewJob(url string, withOption ...func(*Option)) (*Job, error) {

	log.Debug("创建下载任务：%s", url)

	Url, err := netUrl.Parse(url)
	if err != nil {
		return nil, err
	}

	d.lastJobId++

	opt := d.Option.cloneOption()

	for _, set := range withOption {
		set(&opt)
	}

	job := &Job{
		downloader: d,
		Option:     opt,

		id:             d.lastJobId,
		Url:            Url,
		FileName:       filepath.Base(Url.Path),
		FileSize:       0,
		isSupportRange: false,
		isFtp:          Url.Scheme == "ftp",
	}

	job.httpClient = job.newHttpClient()

	return job, nil
}

func (job *Job) TargetFile() string {

	if filepath.IsAbs(job.FileName) {
		return job.FileName
	}

	return filepath.Join(job.Dir, job.FileNamePrefix+job.FileName)
}

func (job *Job) createTargetFile() (*os.File, error) {

	if job.Overwrite {
		return os.Create(job.TargetFile())
	}

	original := job.TargetFile()
	// 检查文件是否存在
	if _, err := os.Stat(original); os.IsNotExist(err) {
		// 文件不存在，返回新建文件
		return os.Create(original)
	}

	index := 1

	base := job.FileNamePrefix + job.FileName
	ext := filepath.Ext(base)
	baseWithoutExt := base[:len(base)-len(ext)]

	for {
		// 构造新的文件名
		newBase := fmt.Sprintf("%s(%d)%s", baseWithoutExt, index, ext)
		newPath := filepath.Join(job.Dir, newBase)

		// 检查文件是否存在
		if _, err := os.Stat(newPath); os.IsNotExist(err) {

			job.FileName = strings.TrimPrefix(newBase, job.FileNamePrefix)

			// 文件不存在，返回新建文件
			return os.Create(job.TargetFile())
		}

		// 文件存在，增加索引并重试
		index++
	}

}

func (job *Job) AvaiableTreads() int {

	// 不支持多线程下载
	if !job.isSupportRange {
		return 1
	}

	if job.FileSize < job.MinChunkSize {
		return 1
	}

	// 如果最小分块大小错误，则单线程下载
	if job.MinChunkSize <= 0 {
		return 1
	}

	threads := int(math.Ceil(float64(job.FileSize) / float64(job.MinChunkSize)))

	if threads > job.MaxThreads {
		return job.MaxThreads
	}

	if threads < 1 {
		return 1
	}

	return threads
}

func (job *Job) Download() (string, error) {

	// 下载之前的拦截器
	job.Interceptors.BeforeDownload(job)

	var err error
	if job.isFtp {
		job.logDebug("下载FTP文件")
		err = job.downloadFtp()
	} else {
		job.logDebug("下载HTTP文件")
		err = job.downloadHttp()
	}

	if err != nil {
		return "", err
	}

	targetFile := job.TargetFile()

	job.logDebug("下载完成， 文件路径：%s", targetFile)

	return targetFile, nil
}

func (job *Job) downloadFtp() error {

	ftpHost := job.Url.Host

	if job.Url.Port() == "" {
		ftpHost += ":21"
	}

	c, err := ftp.Dial(ftpHost, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		newErr := errors.New("FTP 链接出错：" + err.Error())
		job.logErr(newErr.Error())
		return newErr
	}

	job.logDebug("链接 %s 成功，准备登录...", ftpHost)

	username := job.Url.User.Username()
	password, _ := job.Url.User.Password()

	if username == "" {
		username = "anonymous"
	}

	err = c.Login(username, password)
	defer c.Quit()
	if err != nil {
		newErr := errors.New("FTP 登录出错：" + err.Error())
		job.logErr(newErr.Error())
		return newErr
	}

	job.logDebug("登录 %s 成功，准备下载文件...", username)

	// 保存文件之前的拦截器
	job.Interceptors.BeforeSaveFile(job)

	if job.EnableSaveFileDialog {
		job.logDebug("准备打开文件选择对话框... TargetFile() %s", job.TargetFile())
		if path, ok := openSaveFileDialog(job.TargetFile()); ok {
			dir, file := filepath.Split(path)
			job.FileName = file
			job.Dir = dir
			job.Overwrite = true
		} else {
			job.logDebug("用户取消保存。")
			return nil
		}
	}

	if job.FileName == "" {
		return errors.New("文件名不能为空！")
	}

	job.logDebug("创建任务 %s", job.Url.String())

	r, err := c.Retr(job.Url.Path)
	if err != nil {
		return err
	}
	defer r.Close()

	// 打开文件准备写入
	file, err := job.createTargetFile()
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := io.ReadAll(job.Interceptors.FtpDownloading(r, *job.Url))
	if err != nil {
		return err
	}

	_, err = file.Write(buf)

	return err
}

func (job *Job) newHttpClient() *http.Client {

	// 创建一个cookie jar
	jar, err := cookiejar.New(nil)
	if err == nil {
		for _, cookie := range job.Cookies {
			// dat文件里存在多个域名的cookie，所以需要判断域名是否匹配
			if strings.Contains(job.Url.Host, cookie.Domain) {
				jar.SetCookies(job.Url, []*http.Cookie{cookie})
			}
		}
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			// 忽略证书验证
			InsecureSkipVerify: job.InsecureSkipVerify,
		},
	}

	client := &http.Client{
		Transport: tr,
		Jar:       jar,
		Timeout:   job.Timeout,
	}

	return client
}

func (job *Job) downloadHttp() error {

	// 301 跳转会导致 head 失败
	// 不管是否失败，都进行下载（将会使用单线程模式下载）
	job.fetchInfo()

	job.logDebug("创建任务 %s", job.Url)

	// 保存文件之前的拦截器
	job.Interceptors.BeforeSaveFile(job)

	if job.EnableSaveFileDialog {

		if path, ok := openSaveFileDialog(job.TargetFile()); ok {
			dir, file := filepath.Split(path)
			job.FileName = file
			job.Dir = dir
			job.Overwrite = true
		} else {
			job.logDebug("用户取消保存。")
			return errors.New("用户取消保存。")
		}
	}

	if job.FileName == "" || !strings.Contains(job.FileName, ".") {
		job.logDebug("文件名不正确: %s", job.FileName)
		return errors.New("文件名不正确。")
	}

	if job.isSupportRange {

		if err := job.multiThreadDownload(); err != nil {
			job.logErr(err.Error())
			return err
		}
	} else {
		if err := job.singleThreadDownload(); err != nil {
			job.logErr(err.Error())
			return err
		}
	}

	job.logDebug("下载完成：%s", job.TargetFile())
	return nil
}

func (job *Job) singleThreadDownload() error {

	job.logDebug("文件将以单进程模式下载。")

	// 打开文件准备写入
	file, err := job.createTargetFile()
	if err != nil {
		return err
	}
	defer file.Close()

	// 实现默认下载逻辑
	req, err := http.NewRequest(http.MethodGet, job.Url.String(), nil)
	if err != nil {
		return err
	}
	r, err := job.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	reader := job.Interceptors.HttpDownloading(r, req)
	_, err = io.Copy(file, reader)
	return err
}

func (job *Job) multiThreadDownload() error {

	theads := job.AvaiableTreads()
	job.logDebug("文件将以多线程进行下载，线程：%d", theads)

	// 打开文件准备写入
	file, err := job.createTargetFile()
	if err != nil {
		return err
	}
	defer file.Close()

	var wg sync.WaitGroup
	var ctx, cancel = context.WithCancel(context.Background())
	var errs []error
	var errLock sync.Mutex

	defer cancel() // 取消所有goroutine

	// 计算每个线程的分块大小
	chunkSize := uint64(math.Ceil(float64(job.FileSize) / float64(theads)))

	for i := 0; i < theads; i++ {
		start := uint64(i) * chunkSize
		end := start + chunkSize - 1

		// 如果是最后一个部分，加上余数
		if i == theads-1 {
			end = job.FileSize - 1
		}

		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			retry := 0
			for {
				select {
				case <-ctx.Done():
					// 如果收到取消信号，直接返回
					return
				default:
					// 尝试下载分块
					err := job.downloadChunk(file, start, end)

					if err == nil {
						job.logDebug("切片 %d 下载完成", index+1)
						return
					}

					// 如果重试超过3次，记录错误并触发取消操作
					if retry >= 3 {
						errLock.Lock()
						errs = append(errs, err)
						errLock.Unlock()
						cancel() // 取消所有goroutine
						return
					}

					retry++
				}
			}
		}(i)
	}

	wg.Wait() // 等待所有goroutine完成

	if len(errs) > 0 {
		return errs[0] // 返回第一个遇到的错误
	}

	return nil
}

// downloadChunk 下载文件的单个分块
func (job *Job) downloadChunk(file *os.File, start, end uint64) error {

	req, err := http.NewRequest("GET", job.Url.String(), nil)
	if err != nil {
		return err
	}

	// 设置Range头实现断点续传
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := job.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查服务器是否支持Range请求
	if resp.StatusCode != http.StatusPartialContent {
		return errors.New("server doesn't support Range requests")
	}

	// 锁定互斥锁以安全地写入文件
	job._lck.Lock()
	defer job._lck.Unlock()

	// 写入文件的当前位置
	if _, err = file.Seek(int64(start), io.SeekStart); err != nil {
		return err
	}

	reader := job.Interceptors.HttpDownloading(resp, req)

	// 将HTTP响应的Body内容写入到文件中
	_, err = io.Copy(file, reader)
	return err
}

func (job *Job) fetchInfo() {

	req, err := http.NewRequest("HEAD", job.Url.String(), nil)
	if err != nil {
		job.logErr("创建 fetchInfo HEAD 请求失败： %s", err.Error())
		return
	}

	r, err := job.httpClient.Do(req)
	if err != nil {
		job.logErr("fetchInfo 失败： %s", err.Error())
		return
	}
	defer r.Body.Close()

	if r.StatusCode == 404 {
		job.logErr("文件不存在(404)： %s")
		return
	}

	if r.StatusCode > 299 {
		job.logErr("连接出错(%d)： %s", r.StatusCode, job.Url.String())
		return
	}

	// 获取文件名
	job.FileName = getFileNameByResponse(r)

	// 检查是否支持 断点续传
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
	if r.Header.Get("Accept-Ranges") == "bytes" {
		job.isSupportRange = true
	}

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Length
	// 获取文件总大小
	contentLength, err := strconv.ParseUint(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		job.isSupportRange = false
		job.FileSize = 0
		return
	}

	job.FileSize = contentLength
}

func getFileNameByResponse(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)

		if err == nil {
			for k, v := range params {
				// 忽略大小写
				if strings.ToLower(k) == "filename" {
					return v
				}
			}
		}

	}
	return filepath.Base(resp.Request.URL.Path)
}

func openSaveFileDialog(filePath string) (filepath string, ok bool) {
	var ofn win.OPENFILENAME
	buf := make([]uint16, syscall.MAX_PATH) // 假设路径可能更长，增加缓冲区大小
	ofn.LStructSize = uint32(unsafe.Sizeof(ofn))
	ofn.LpstrFile = &buf[0]
	ofn.NMaxFile = uint32(len(buf))
	ofn.Flags = win.OFN_OVERWRITEPROMPT

	// UTF16FromString 不支持中间带 \0 的字符串，所以需要手动拼接
	filter, _ := syscall.UTF16FromString("所有文件（*.*）")
	filterM, _ := syscall.UTF16FromString("*.*")
	filter = append(filter, filterM...)
	filter = append(filter, 0)
	ofn.LpstrFilter = &filter[0]

	// 转换文件名到UTF-16，并检查错误
	if utf16FileName, err := syscall.UTF16FromString(filePath); err == nil {
		copy(buf, utf16FileName)
	}

	ok = win.GetSaveFileName(&ofn)

	if ok {
		filepath = syscall.UTF16ToString(buf)
	}

	return
}

func (job *Job) logDebug(tpl string, vars ...interface{}) {
	log.Debug(fmt.Sprintf("[下载任务 %d ]: ", job.id)+tpl, vars...)
}

func (job *Job) logErr(tpl string, vars ...interface{}) {

	log.Error(fmt.Sprintf("[下载任务 %d ]: ", job.id)+tpl, vars...)
}
