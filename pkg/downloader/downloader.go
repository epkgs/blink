package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
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

type Downloader struct {
	lastJobId uint64
	Option

	afterCreateJobInterceptor AfterCreateJobInterceptor
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
	ticker         *time.Ticker

	_lck *sync.Mutex // 用于确保写入文件的顺序
}

type Option struct {
	Dir                  string        // 下载路径，如果为空则使用当前目录
	FileNamePrefix       string        // 文件名前缀，默认空
	MaxThreads           int           // 下载线程，默认4
	MinChunkSize         uint64        // 最小分块大小，默认500KB
	EnableSaveFileDialog bool          // 是否打开保存文件对话框，默认false
	Overwrite            bool          // 是否覆盖已存在的文件，默认false
	Timeout              time.Duration // 超时时间，默认10秒
}

type AfterCreateJobInterceptor func(job *Job)

func (opt Option) cloneOption() Option {
	return Option{
		Dir:                  opt.Dir,
		FileNamePrefix:       opt.FileNamePrefix,
		MaxThreads:           opt.MaxThreads,
		MinChunkSize:         opt.MinChunkSize,
		EnableSaveFileDialog: opt.EnableSaveFileDialog,
		Overwrite:            opt.Overwrite,
		Timeout:              opt.Timeout,
	}
}

func New(withOption ...func(*Option)) *Downloader {

	pwd, err := os.Getwd()
	if err != nil {
		pwd = ""
	}

	// 默认参数
	opt := Option{
		Dir:                  pwd,
		FileNamePrefix:       "",
		MaxThreads:           4,
		MinChunkSize:         500 * 1024, // 500KB
		EnableSaveFileDialog: false,
		Overwrite:            false,
		Timeout:              10 * time.Second,
	}

	for _, set := range withOption {
		set(&opt)
	}

	downloader := &Downloader{
		lastJobId: 0,
		Option:    opt,
	}

	// 空实现
	downloader.afterCreateJobInterceptor = func(job *Job) {}

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
		ticker:         time.NewTicker(opt.Timeout),
	}

	d.afterCreateJobInterceptor(job)

	return job, nil
}

func (d *Downloader) AfterCreateJob(interceptor AfterCreateJobInterceptor) {
	d.afterCreateJobInterceptor = interceptor
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
	select {
	case <-job.ticker.C:
		return "", errors.New("下载超时")
	default:
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

	if job.EnableSaveFileDialog {
		job.logDebug("准备打开文件选择对话框... TargetFile() %s", job.TargetFile())
		if path, ok := openSaveFileDialog(job.TargetFile()); ok {
			dir, file := filepath.Split(path)
			job.FileName = file
			job.Dir = dir
		} else {
			job.logDebug("用户取消保存。")
			return nil
		}
	}

	if job.FileName == "" || !strings.Contains(job.FileName, ".") {
		return fmt.Errorf("文件名不正确: %s", job.FileName)
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

	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	_, err = file.Write(buf)

	return err
}

func (job *Job) downloadHttp() error {

	if err := job.fetchInfo(); err != nil {
		job.logErr("获取文件信息出错：" + err.Error())
		return err
	}

	job.logDebug("创建任务 %s", job.Url)
	if job.EnableSaveFileDialog {

		if path, ok := openSaveFileDialog(job.TargetFile()); ok {
			dir, file := filepath.Split(path)
			job.FileName = file
			job.Dir = dir
		} else {
			job.logDebug("用户取消保存。")
			return nil
		}
	}

	if job.FileName == "" || !strings.Contains(job.FileName, ".") {
		job.logDebug("文件名不正确: %s", job.FileName)
		return errors.New("文件名不正确。")
	}

	if err := job.multiThreadDownload(); err != nil {
		job.logErr(err.Error())
		return err
	}

	job.logDebug("下载完成：%s", job.TargetFile())
	return nil
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
	resp, err := http.DefaultClient.Do(req)
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

	// 将HTTP响应的Body内容写入到文件中
	_, err = io.Copy(file, resp.Body)
	return err
}

func (job *Job) fetchInfo() error {

	r, err := http.Head(job.Url.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode == 404 {
		return fmt.Errorf("文件不存在： %s", job.Url.String())
	}

	if r.StatusCode > 299 {
		return fmt.Errorf("连接 %s 出错。", job.Url.String())
	}

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
		return nil
	}

	job.FileSize = contentLength

	if !job.EnableSaveFileDialog {
		job.FileName = getFileNameByResponse(r)
	}

	return nil
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
