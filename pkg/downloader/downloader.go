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
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/epkgs/blink/internal/log"
	"github.com/epkgs/blink/pkg/alert"
	"github.com/jlaffaye/ftp"
	"github.com/lxn/win"
)

type IDownloadChunkCallback func(res *http.Response, index uint64) error

type IBeforeDownloadInterceptor func(job *Job)
type IHttpDownloadingInterceptor func(job *Job, res *http.Response) io.Reader
type IFtpDownloadingInterceptor func(job *Job, res *ftp.Response) io.Reader
type IBeforeSaveFileInterceptor func(job *Job)

type IInterceptors struct {
	BeforeDownload  IBeforeDownloadInterceptor
	HttpDownloading IHttpDownloadingInterceptor
	FtpDownloading  IFtpDownloadingInterceptor
	BeforeSaveFile  IBeforeSaveFileInterceptor
}

type Config struct {
	Dir            string // 下载路径，如果为空则使用当前目录
	FileNamePrefix string // 文件名前缀，默认空
	MaxThreads     uint64 // 下载线程，默认5
	MinChunkSize   uint64 // 最小分块大小，默认500KB

	Timeout time.Duration  // 超时时间，默认10秒
	Cookies []*http.Cookie // 请求头Cookie，默认空。

	EnableSaveFileDialog bool // 是否打开保存文件对话框，默认false
	OverwriteFile        bool // 是否覆盖已存在的文件，默认false
	InsecureSkipVerify   bool // 跳过证书验证，默认false

	Interceptors IInterceptors // 拦截器
}

func (conf Config) Clone() Config {
	// 在GO中，直接返回是值的浅拷贝，基本类型会直接拷贝，引用类型会拷贝指针
	//
	// Cookies 是 slice，如果需要深拷贝需要手动拷贝
	// 参考：https://stackoverflow.com/questions/32167098/how-to-deep-copy-a-slice
	//
	// Interceptors 字段是一个 IInterceptors 类型的结构体，
	// 但它的所有字段都是函数类型。函数类型在Go中不是引用类型，
	// 它们不存储状态，因此不需要进行深拷贝。
	return conf
}

type Downloader struct {
	Config

	lastJobId uint64
	ctx       context.Context
}

type Job struct {
	Config

	downloader *Downloader

	id             uint64
	Url            *netUrl.URL
	FileName       string
	FileSize       uint64
	isSupportRange bool
	isFtp          bool

	ctx    context.Context
	cancel context.CancelFunc
}

func New(withConfig ...func(*Config)) *Downloader {
	return NewWithContext(context.Background(), withConfig...)
}

func NewWithContext(ctx context.Context, withConfig ...func(*Config)) *Downloader {

	// 默认参数
	defaultOption := Config{
		Dir:            "",
		FileNamePrefix: "",
		MaxThreads:     5,
		MinChunkSize:   500 * 1024, // 500KB

		Timeout: 10 * time.Second,
		Cookies: make([]*http.Cookie, 0),

		EnableSaveFileDialog: false,
		OverwriteFile:        false,
		InsecureSkipVerify:   false,

		Interceptors: IInterceptors{
			BeforeDownload: func(job *Job) {}, // 默认空实现
			HttpDownloading: func(job *Job, res *http.Response) io.Reader {
				return res.Body
			},
			FtpDownloading: func(job *Job, res *ftp.Response) io.Reader {
				return res
			},
			BeforeSaveFile: func(job *Job) {}, // 默认空实现
		},
	}

	downloader := &Downloader{
		lastJobId: 0,
		Config:    defaultOption,

		ctx: ctx,
	}

	return downloader.WithConfig(withConfig...)
}

// 修改 Downloader 的默认参数
func (d *Downloader) WithConfig(withConfig ...func(*Config)) *Downloader {
	for _, set := range withConfig {
		set(&d.Config)
	}

	return d
}

func (d *Downloader) WithContext(ctx context.Context) *Downloader {
	d.ctx = ctx
	return d
}

func (d *Downloader) Download(url string, withConfig ...func(*Config)) (targetFile string, err error) {
	job, err := d.newJob(url, withConfig...)
	if err != nil {
		return "", err
	}
	return job.download()
}

func (d *Downloader) newJob(url string, withConfig ...func(*Config)) (*Job, error) {

	Url, err := netUrl.Parse(url)
	if err != nil {
		return nil, err
	}

	d.lastJobId++

	conf := d.Config.Clone()

	for _, set := range withConfig {
		set(&conf)
	}

	job := &Job{
		downloader: d,
		Config:     conf,

		id:             d.lastJobId,
		Url:            Url,
		FileName:       filepath.Base(Url.Path),
		FileSize:       0,
		isSupportRange: false,
		isFtp:          Url.Scheme == "ftp",
	}

	job.ctx, job.cancel = context.WithCancel(d.ctx)

	return job, nil
}

func (job *Job) targetFile() string {

	if filepath.IsAbs(job.FileName) {
		return job.FileName
	}

	return filepath.Join(job.Dir, job.FileNamePrefix+job.FileName)
}

func (job *Job) handleSaveFileDialog() error {

	if job.EnableSaveFileDialog {

		fileNameOk := false

		for {
			if fileNameOk {
				break
			}

			path, ok := openSaveFileDialog(job.targetFile())
			if !ok {
				return errors.New("用户取消保存。")
			}

			dir, fname := filepath.Split(path)

			if strings.TrimSpace(fname) == "" {
				alert.Error("文件名不能为空")
				continue
			}

			job.FileName = fname
			job.Dir = dir
			job.OverwriteFile = true

			fileNameOk = true // 文件名正确，跳出循环
		}

	}

	return nil
}

func (job *Job) getFinalTargetFile() string {

	if job.OverwriteFile {
		return job.targetFile()
	}

	original := job.targetFile()
	// 检查文件是否存在
	if _, err := os.Stat(original); os.IsNotExist(err) {
		// 文件不存在，返回新建文件
		return original
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
			return job.targetFile()
		}

		// 文件存在，增加索引并重试
		index++
	}

}

func avaiableTreads(fileSize, minChunkSize, maxThreads uint64) uint64 {

	if fileSize < minChunkSize {
		return 1
	}

	// 如果最小分块大小错误，则单线程下载
	if minChunkSize <= 0 {
		return 1
	}

	threads := uint64(math.Ceil(float64(fileSize) / float64(minChunkSize)))

	if threads > maxThreads {
		return maxThreads
	}

	if threads < 1 {
		return 1
	}

	return threads
}

func (job *Job) download() (targetFile string, err error) {

	defer job.cancel()

	// 下载之前的拦截器
	job.Interceptors.BeforeDownload(job)

	var tmpFiles []string

	defer func() {
		for _, f := range tmpFiles {
			os.Remove(f)
		}
	}()

	if job.isFtp {
		job.logDebug("创建FTP下载任务：%s", job.Url.String())
		tmpFiles, err = job.downloadFtp()
	} else {
		job.logDebug("创建HTTP下载任务：%s", job.Url.String())
		tmpFiles, err = job.downloadHttp()
	}

	// 等待下载完成
	if err != nil {
		return "", err
	}

	targetFile = job.getFinalTargetFile()

	err = mergeFiles(tmpFiles, targetFile)
	if err != nil {
		os.Remove(targetFile)
		job.logErr("将临时文件写入目标文件失败：%s", err.Error())
		return "", err
	}

	if err == nil {
		job.logDebug("下载完成， 文件路径：%s", targetFile)
	}

	return targetFile, err
}

func (job *Job) downloadFtp() (tmpFiles []string, err error) {

	defer func() {
		if err != nil {
			// 下载失败，清理临时文件
			for _, f := range tmpFiles {
				os.Remove(f)
			}
			tmpFiles = nil
		}
	}()

	ftpHost := job.Url.Host

	if job.Url.Port() == "" {
		ftpHost += ":21"
	}

	tmpFiles = make([]string, 0)

	c, err := ftp.Dial(ftpHost, ftp.DialWithContext(job.ctx), ftp.DialWithTimeout(job.Timeout))
	if err != nil {
		newErr := errors.New("FTP 链接出错：" + err.Error())
		job.logErr(newErr.Error())
		return tmpFiles, newErr
	}

	job.logDebug("链接 %s 成功，准备登录...", ftpHost)

	username := job.Url.User.Username()
	password, _ := job.Url.User.Password()

	if username == "" {
		username = "anonymous"
	}

	err = c.Login(username, password)
	defer func() {
		_ = c.Quit()
	}()

	if err != nil {
		newErr := errors.New("FTP 登录出错：" + err.Error())
		job.logErr(newErr.Error())
		return tmpFiles, newErr
	}

	job.logDebug("登录 %s 成功，准备下载文件...", username)

	res, err := c.Retr(job.Url.Path)
	if err != nil {
		return tmpFiles, err
	}
	defer res.Close()

	// 准备临时文件
	file, err := job.getTempFile()
	if err != nil {
		return tmpFiles, err
	}
	defer file.Close()
	tmpFiles = append(tmpFiles, file.Name())

	buf, err := io.ReadAll(job.Interceptors.FtpDownloading(job, res))
	if err != nil {
		return tmpFiles, err
	}

	_, err = file.Write(buf)

	return tmpFiles, nil
}

func (job *Job) sentRequest(req *http.Request) (*http.Response, error) {

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
		// Timeout:   job.Timeout,
	}

	return client.Do(req)
}

func (job *Job) getTempFile() (*os.File, error) {

	dir := path.Join(os.TempDir(), "mini-blink")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return os.CreateTemp(dir, "download_*.tmp")
}

func (job *Job) getInfoByResponse(res *http.Response) {

	// 获取文件名
	job.FileName = getFileNameByResponse(res)

	// 检查是否支持 断点续传
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
	if res.Header.Get("Accept-Ranges") == "bytes" {
		job.isSupportRange = true
	}

	// 通过 Content-Range 获取文件大小
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Range
	contentRange := res.Header.Get("Content-Range")
	if contentRange != "" {
		crs := strings.Split(contentRange, "/")
		if len(crs) == 2 && crs[1] != "*" {

			if size, err := strconv.ParseUint(crs[1], 10, 64); err == nil {
				job.FileSize = size
			}
		}
	} else {

		// 当未设置 Content-Rnage 时，使用 Content-Length 获取文件总大小
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Length
		if contentLength, err := strconv.ParseUint(res.Header.Get("Content-Length"), 10, 64); err == nil {
			job.FileSize = contentLength
		}
	}
}

// 多线程下载。返回下载后的临时文件和错误
func (job *Job) downloadHttp() (tmpFiles []string, downloadErr error) {

	defer func() {
		if downloadErr != nil {
			job.logErr(downloadErr.Error())
			// 下载失败，清理临时文件
			for _, f := range tmpFiles {
				os.Remove(f)
			}
			tmpFiles = nil
		}
	}()

	chunkStart := uint64(0)
	chunkEnd := job.MinChunkSize - 1
	tmpFiles = make([]string, 1) // 预初始化为1个元素的数组

	ctx, cancel := context.WithCancel(job.ctx)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		// 尝试用多线程下载的方式，以最小切片大小下载第一部分
		fname, err := job.downloadChunk(ctx, 0, fmt.Sprintf("%d-%d", chunkStart, chunkEnd), func(res *http.Response, index uint64) error {

			// 获取文件信息
			job.getInfoByResponse(res)

			wg.Add(1)
			go func() {
				defer wg.Done()

				// 使用协程处理另存为选择框，使其不阻塞下载

				// 保存文件之前的拦截器
				job.Interceptors.BeforeSaveFile(job)

				if err := job.handleSaveFileDialog(); err != nil {
					downloadErr = fmt.Errorf("保存文件失败：%s", err.Error())
					cancel()
				}

			}()

			// 如果返回的状态码为 206，则表示服务器支持断点续传，需要继续下载剩余部分
			if res.StatusCode == http.StatusPartialContent {
				if job.FileSize == 0 {

					// 支持断点续传，但无法获取到文件大小，则再加一个线程下载完剩余的部分
					job.logDebug("服务器支持断点续传，但无法获取文件大小，将以新进程继续下载剩余部分")

					tmpFiles = make([]string, 2)
					wg.Add(1)
					go func() {
						defer wg.Done()
						fname2, err := job.downloadChunk(ctx, 1, fmt.Sprintf("%d-", chunkEnd+1))
						if err != nil {
							downloadErr = fmt.Errorf("下载失败：%s", err.Error())
							cancel()
							return
						}
						tmpFiles[1] = fname2
					}()
				} else {
					// 支持断点续传，且获取到文件大小，计算可用线程数
					remainSize := job.FileSize - chunkEnd - 1 // 剩余大小
					remainTheads := avaiableTreads(
						remainSize, // 剩余大小
						job.MinChunkSize,
						job.MaxThreads-1, // 扣除已使用的线程数
					)
					theads := remainTheads + 1

					job.logDebug("服务器支持断点续传，文件将以多线程继续下载，线程：%d，文件大小：%d", theads, job.FileSize)

					tmpFiles = make([]string, theads)

					// 循环发起剩余下载请求部分，使用协程，不阻塞第一个下载进程
					wg.Add(1)
					go func() {
						defer wg.Done()

						// 计算每个线程的分块大小
						chunkSize := uint64(math.Ceil(float64(remainSize) / float64(remainTheads)))

						for idx := uint64(1); idx < theads; idx++ {
							chunkStart = chunkEnd + 1
							chunkEnd = chunkStart + chunkSize - 1

							// 如果是最后一个部分，end 为文件末端
							if idx == theads-1 {
								chunkEnd = job.FileSize - 1
							}

							wg.Add(1)
							go func(index uint64) {
								defer wg.Done()

								retry := 0
								for {
									select {
									case <-ctx.Done():
										// 如果收到取消信号，直接返回
										return
									default:
										retry++

										// 尝试下载分块
										fname, err := job.downloadChunk(ctx, index, fmt.Sprintf("%d-%d", chunkStart, chunkEnd))
										if err == nil {
											tmpFiles[index] = fname
											return
										}

										// 如果重试超过3次，记录错误并触发取消操作
										if retry >= 3 {
											downloadErr = err
											cancel()
											return
										}
									}
								}
							}(idx)
						}

					}()

				}
			}

			return nil
		})

		tmpFiles[0] = fname // 等待结束再赋值，因为 callback 内部会重新初始化 tmpFiles

		if err != nil {
			downloadErr = err
			cancel()
		}

	}()

	wg.Wait() // 等待所有下载子线程完成

	return tmpFiles, downloadErr
}

// 下载文件的单个分块，带回调函数
func (job *Job) downloadChunk(ctx context.Context, index uint64, byteRange string, callbacks ...IDownloadChunkCallback) (tmpFile string, err error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.Url.String(), nil)
	if err != nil {
		return "", err
	}

	// 设置Range头实现断点续传
	req.Header.Set("Range", fmt.Sprintf("bytes=%s", byteRange))

	res, err := job.sentRequest(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode != http.StatusRequestedRangeNotSatisfiable {
		return "", fmt.Errorf("下载失败，服务器返回错误状态码：%d", res.StatusCode)
	}

	for _, callback := range callbacks {

		if err := callback(res, index); err != nil {
			return "", err
		}
	}

	reader := job.Interceptors.HttpDownloading(job, res)

	// 准备临时文件
	file, err := job.getTempFile()
	if err != nil {
		return "", err
	}
	defer file.Close()
	tmpFile = file.Name()

	job.logDebug("[ 线程 %d ] 下载到临时文件 %s", index+1, tmpFile)

	// 将HTTP响应的Body内容写入到文件中
	_, err = io.Copy(file, reader)
	return tmpFile, err
}

func getFileNameByResponse(res *http.Response) string {
	contentDisposition := res.Header.Get("Content-Disposition")
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
	return filepath.Base(res.Request.URL.Path)
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

// mergeFiles 实现跨卷/分区移动文件
func mergeFiles(sourcePaths []string, destPath string) error {

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	for _, sourcePath := range sourcePaths {

		err := func() error {
			inputFile, err := os.Open(sourcePath)
			if err != nil {
				return fmt.Errorf("Couldn't open source file: %s", err)
			}
			defer inputFile.Close()

			_, err = io.Copy(outputFile, inputFile)
			if err != nil {
				return fmt.Errorf("Writing to output file failed: %s", err)
			}
			return nil
		}()

		if err != nil {
			return err
		}
	}
	return nil
}
