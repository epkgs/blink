package main

import (
	"compress/flate"
	"io"
	"net/http"

	"github.com/epkgs/blink"
	"github.com/epkgs/blink/pkg/downloader"
)

func main() {
	app := blink.NewApp()
	defer app.Exit()

	// 单线程直接下载
	_, _ = app.Download("https://httpbin.org/robots.txt")

	// deflate压缩内容下载
	_, _ = app.Download("https://comment.bilibili.com/177987845.xml", func(o *downloader.Option) {
		o.Interceptors.HttpDownloading = func(job *downloader.Job, res *http.Response) io.Reader {
			// 检查Content-Encoding是否为deflate
			contentEncoding := res.Header.Get("Content-Encoding")
			if contentEncoding == "deflate" {
				// 如果是deflate编码，解压缩数据
				return flate.NewReader(res.Body)
			} else {
				// 如果不是deflate编码，直接将响应体内容写入文件
				return res.Body
			}
		}
	})

	// 保存文件之前修改文件名
	_, _ = app.Download("https://httpbin.org/image", func(o *downloader.Option) {
		o.Interceptors.BeforeSaveFile = func(job *downloader.Job) {
			// 修改文件名
			job.FileName = "test.png"
		}
	})

	// 多线程下载
	_, _ = app.Download("https://github.com/iBotPeaches/Apktool/releases/download/v2.9.3/apktool_2.9.3.jar")
}
