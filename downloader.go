package blink

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"net/url"
	netUrl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/epkgs/mini-blink/internal/log"
	"github.com/jlaffaye/ftp"
	"github.com/lxn/win"
)

type Downloader struct {
	lastJobId uint64
	DownloadOption
}

type DownloadJob struct {
	downloader *Downloader
	DownloadOption

	id             uint64
	url            *netUrl.URL
	dir            string
	filename       string
	size           int64
	isSupportRange bool
	isFtp          bool
}

type DownloadOption struct {
	Dir                  string // 下载路径，如果为空则使用当前目录
	Threads              int    // 下载线程
	EnableSaveFileDialog bool   // 是否打开保存文件对话框
}

func (opt DownloadOption) cloneOption() DownloadOption {
	return DownloadOption{
		Dir:                  opt.Dir,
		Threads:              opt.Threads,
		EnableSaveFileDialog: opt.EnableSaveFileDialog,
	}
}

func NewDownloader(withOption ...func(*DownloadOption)) *Downloader {

	pwd, err := os.Getwd()
	if err != nil {
		pwd = ""
	}

	// 默认参数
	opt := DownloadOption{
		Dir:                  pwd,
		Threads:              4,
		EnableSaveFileDialog: true,
	}

	for _, set := range withOption {
		set(&opt)
	}

	downloader := &Downloader{
		lastJobId:      0,
		DownloadOption: opt,
	}

	return downloader
}

func (d *Downloader) Download(url string, withOption ...func(*DownloadOption)) error {
	job, err := d.NewJob(url, withOption...)
	if err != nil {
		return err
	}
	return job.Download()
}

func (d *Downloader) NewJob(url string, withOption ...func(*DownloadOption)) (*DownloadJob, error) {

	Url, err := netUrl.Parse(url)
	if err != nil {
		return nil, err
	}

	d.lastJobId++

	opt := d.DownloadOption.cloneOption()

	for _, set := range withOption {
		set(&opt)
	}

	job := &DownloadJob{
		downloader:     d,
		DownloadOption: opt,

		id:             d.lastJobId,
		url:            Url,
		filename:       filepath.Base(Url.Path),
		size:           0,
		isSupportRange: false,
		isFtp:          Url.Scheme == "ftp",
	}

	return job, nil
}

func (job *DownloadJob) TargetFile() string {
	return filepath.Join(job.Dir, job.filename)
}

func (job *DownloadJob) Download() error {

	if job.isFtp {
		return job.downloadFtp()
	}

	return job.downloadHttp()
}

func (job *DownloadJob) downloadFtp() error {

	c, err := ftp.Dial(job.url.Host+job.url.Port(), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		job.logErr("打开 FTP 出错：" + err.Error())
		return errors.New("链接 FTP 服务器出错")
	}

	username := job.url.User.Username()
	password, _ := job.url.User.Password()

	if username == "" {
		username = "anonymous"
	}

	err = c.Login(username, password)
	defer c.Quit()
	if err != nil {
		job.logErr("登录 FTP 出错：" + err.Error())
		return errors.New("登录 FTP 出错")
	}

	if job.EnableSaveFileDialog {
		if path, ok := openSaveFileDialog(job.TargetFile()); ok {
			dir, file := filepath.Split(path)
			job.filename = file
			job.Dir = dir
		} else {
			job.logDebug("用户取消保存。")
			return nil
		}
	}

	job.logDebug("创建任务 %s", job.url.String())

	r, err := c.Retr(job.url.Path)
	if err != nil {
		return err
	}
	defer r.Close()

	// 打开文件准备写入
	file, err := os.Create(job.TargetFile())
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

func (job *DownloadJob) downloadHttp() error {

	if err := job.fetchInfo(); err != nil {
		job.logErr("获取文件信息出错：" + err.Error())
		return err
	}

	job.logDebug("创建任务 %s", job.url)
	if job.EnableSaveFileDialog {

		if path, ok := openSaveFileDialog(job.TargetFile()); ok {
			dir, file := filepath.Split(path)
			job.filename = file
			job.Dir = dir
		} else {
			job.logDebug("用户取消保存。")
			return nil
		}
	}

	if job.isSupportRange {
		job.logDebug("支持断点续传，线程：%d", job.Threads)
		if err := job.multiThreadDownload(); err != nil {
			job.logErr(err.Error())
			return err
		}
	} else {
		job.logDebug("不支持断点续传，将以单进程模式下载。")
		if err := job.singleThreadDownload(); err != nil {
			job.logErr(err.Error())
			return err
		}
	}

	job.logDebug("下载完成：%s", job.TargetFile())
	return nil
}

func (job *DownloadJob) singleThreadDownload() error {

	// 打开文件准备写入
	file, err := os.Create(job.TargetFile())
	if err != nil {
		return err
	}
	defer file.Close()

	// 实现默认下载逻辑
	// 使用http.Get或其他方式下载整个文件
	resp, err := http.Get(job.url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (job *DownloadJob) multiThreadDownload() error {
	// 打开文件准备写入
	file, err := os.Create(job.TargetFile())
	if err != nil {
		return err
	}
	defer file.Close()

	var wg sync.WaitGroup
	var mutex sync.Mutex // 用于确保写入文件的顺序
	var ctx, cancel = context.WithCancel(context.Background())
	var errs []error
	var errLock sync.Mutex

	// 计算每个线程的分块大小
	chunkSize := int64(math.Ceil(float64(job.size) / float64(job.Threads)))

	for i := 0; i < job.Threads; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1

		// 如果是最后一个部分，加上余数
		if i == job.Threads-1 {
			end = job.size - 1
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
					err := job.downloadChunk(&mutex, file, start, end)

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
func (job *DownloadJob) downloadChunk(mutex *sync.Mutex, file *os.File, start, end int64) error {

	req, err := http.NewRequest("GET", job.url.String(), nil)
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
	mutex.Lock()
	defer mutex.Unlock()

	// 写入文件的当前位置
	if _, err = file.Seek(start, io.SeekStart); err != nil {
		return err
	}

	// 将HTTP响应的Body内容写入到文件中
	_, err = io.Copy(file, resp.Body)
	return err
}

func (job *DownloadJob) fetchInfo() error {

	r, err := http.Head(job.url.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode == 404 {
		return fmt.Errorf("文件不存在： %s", job.url.String())
	}

	if r.StatusCode > 299 {
		return fmt.Errorf("连接 %s 出错。", job.url.String())
	}

	// 检查是否支持 断点续传
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
	if r.Header.Get("Accept-Ranges") == "bytes" {
		job.isSupportRange = true
	}

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Length
	// 获取文件总大小
	contentLength, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		job.isSupportRange = false
		job.size = 0
		return nil
	}

	job.size = contentLength

	return nil
}

func getFileNameByResponse(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)

		if err != nil {
			return getFileNameByUrl(resp.Request.URL.Path)
		}
		return params["filename"]
	}
	return getFileNameByUrl(resp.Request.URL.Path)
}

func getFileNameByUrl(downloadUrl string) string {
	parsedUrl, _ := url.Parse(downloadUrl)
	return filepath.Base(parsedUrl.Path)
}

func openSaveFileDialog(fileName string) (filepath string, ok bool) {
	var ofn win.OPENFILENAME
	buf := make([]uint16, syscall.MAX_PATH) // 假设路径可能更长，增加缓冲区大小
	ofn.LStructSize = uint32(unsafe.Sizeof(ofn))
	ofn.LpstrFile = &buf[0]
	ofn.NMaxFile = uint32(len(buf))
	ofn.Flags = win.OFN_OVERWRITEPROMPT

	filter := StringToU16Arr("所有文件（*.*）\000*.*\000\000")
	ofn.LpstrFilter = &filter[0]

	// 转换文件名到UTF-16，并检查错误
	if utf16FileName, err := syscall.UTF16FromString(fileName); err == nil {
		copy(buf, utf16FileName)
	}

	ok = win.GetSaveFileName(&ofn)

	if ok {
		filepath = syscall.UTF16ToString(buf)
	}

	return
}

func (job *DownloadJob) logDebug(tpl string, vars ...interface{}) {
	log.Debug(fmt.Sprintf("[下载任务 %d ]: ", job.id)+tpl, vars...)
}

func (job *DownloadJob) logErr(tpl string, vars ...interface{}) {

	log.Error(fmt.Sprintf("[下载任务 %d ]: ", job.id)+tpl, vars...)
}
