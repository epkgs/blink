package blink

import (
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var mutex sync.Mutex // 用于确保写入文件的顺序

type Downloader struct {
	lastJobId uint64
	threads   int // 下载线程
}

type DownloadJob struct {
	Downloader *Downloader

	id           uint64
	url          string
	target       string
	size         int64
	supportRange bool
}

func NewDownloader(threads int) *Downloader {
	return &Downloader{
		lastJobId: 1,
		threads:   threads,
	}
}

func (d *Downloader) Download(url string) error {
	job := &DownloadJob{

		Downloader: d,

		id:           d.lastJobId,
		url:          url,
		target:       getFileNameByUrl(url),
		size:         0,
		supportRange: false,
	}

	if err := job.fetchInfo(); err != nil {
		job.msgErr("获取文件信息出错：" + err.Error())
		return err
	}

	d.lastJobId++

	job.logInfo("创建任务 %s", url)

	if target, ok := openSaveFileDialog(job.target); ok {
		job.target = target
	} else {
		job.logInfo("用户取消保存。")
		return nil
	}

	if job.supportRange {
		job.logInfo("支持断点续传，线程：%d", job.Downloader.threads)
		if err := job.multiThreadDownload(); err != nil {
			job.msgErr(err.Error())
			return err
		}
	} else {
		job.logInfo("不支持断点续传，将以单进程模式下载。")
		if err := job.singleThreadDownload(); err != nil {
			job.msgErr(err.Error())
			return err
		}
	}

	job.logInfo("下载完成：%s", job.target)
	return nil
}

func (job *DownloadJob) logInfo(tpl string, vars ...any) {
	logInfo(fmt.Sprintf("[下载任务 %d ]: ", job.id)+tpl, vars...)
}

func (job *DownloadJob) msgErr(lines ...string) int32 {
	return MessageBoxError(0, fmt.Sprintf("[下载任务 %d ]: ", job.id), lines...)
}

func (job *DownloadJob) singleThreadDownload() error {
	// 实现默认下载逻辑
	// 使用http.Get或其他方式下载整个文件
	resp, err := http.Get(job.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 打开文件准备写入
	file, err := os.Create(job.target)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (job *DownloadJob) multiThreadDownload() error {

	// 打开文件准备写入
	file, err := os.Create(job.target)
	if err != nil {
		return err
	}
	defer file.Close()

	threads := job.Downloader.threads

	// 计算每个线程的分块大小
	chunkSize := int64(math.Ceil(float64(job.size) / float64(threads)))
	errs := []error{}

	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1

		// 如果是最后一个部分，加上余数
		if i == threads-1 {
			end += job.size - 1
		}

		wg.Add(1)

		go func() {

			if err := job.downloadChunk(file, start, end); err != nil {
				errs = append(errs, err)
			}

			job.logInfo("切片 %d 下载完成", i+1)
			wg.Done()

		}()

	}

	wg.Wait() // 等待所有goroutine完成

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

// downloadChunk 下载文件的单个分块
func (job *DownloadJob) downloadChunk(file *os.File, start, end int64) error {

	req, err := http.NewRequest("GET", job.url, nil)
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

	r, err := http.Head(job.url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// 检查是否支持 断点续传
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
	if r.Header.Get("Accept-Ranges") == "bytes" {
		job.supportRange = true
	}

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Length
	// 获取文件总大小
	contentLength, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
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
	buf := make([]uint16, 260)
	ofn.LStructSize = uint32(unsafe.Sizeof(ofn))
	ofn.LpstrFile = &buf[0]
	ofn.NMaxFile = uint32(len(buf))
	ofn.Flags = win.OFN_OVERWRITEPROMPT

	filter := StringToU16Arr("所有文件（*.*）\000*.*\000\000")
	ofn.LpstrFilter = &filter[0]

	if uints, err := syscall.UTF16FromString(fileName); err == nil {
		copy(buf, uints)
	}

	ok = win.GetSaveFileName(&ofn)

	filepath = syscall.UTF16ToString(buf)

	return

}
