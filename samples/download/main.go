package main

import (
	blink "github.com/epkgs/blink"
	"github.com/epkgs/blink/pkg/alert"
	"github.com/epkgs/blink/pkg/downloader"
	"os"
)

func main() {
	app := blink.NewApp()
	defer app.Exit()

	view := app.CreateWebWindowPopup(blink.WithWebWindowSize(1200, 900))
	view.Window.SetTitle("下载测试")
	view.Window.MoveToCenter()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	// 直接下载
	//app.Downloader.Download("https://httpbin.org/robots.txt") // 空文件
	//app.Downloader.DownloadFile("https://httpbin.org/robots.txt")
	//app.Downloader.DownloadFile("https://httpbin.org/drip?duration=1&numbytes=100&code=200&delay=2")
	//app.Downloader.DownloadFile("https://httpbin.org/stream-bytes/100")

	// deflate压缩内容下载
	//app.Downloader.Download("https://comment.bilibili.com/177987845.xml") // 空文件
	//app.Downloader.DownloadFile("https://comment.bilibili.com/177987845.xml")
	//app.Downloader.DownloadFile("https://httpbin.org/deflate")

	// 图片下载,文件名没有获取到
	//app.Downloader.DownloadFile("https://httpbin.org/image")

	view.LoadURL("https://catalog.ldc.upenn.edu/login")
	view.ShowWindow()
	/*
		用户名: duyu20@mails.tsinghua.edu.cn
		密码: algebra-2023
		右侧有个Downloads连接：https://catalog.ldc.upenn.edu/organization/downloads
	*/
	/*
		miniblink的cookie.dat：
			# Netscape HTTP Cookie File
			# https://curl.haxx.se/docs/http-cookies.html
			# This file was generated by libcurl! Edit at your own risk.

			.upenn.edu	TRUE	/	FALSE	1723686255	_gat	1
			#HttpOnly_.ldc.upenn.edu	TRUE	/	TRUE	0	_xiexie	d425bbf084773bbb8e957b209a603a2a
			.upenn.edu	TRUE	/	FALSE	1786758196	_ga_J9ZV05KG56	GS1.2.1723685408.1.1.1723686196.0.0.0
			.upenn.edu	TRUE	/	FALSE	1723772595	_gid	GA1.2.1524576267.1723685407
			.upenn.edu	TRUE	/	FALSE	1786758195	_ga	GA1.2.1710216791.1723685407
			catalog.ldc.upenn.edu	FALSE	/	TRUE	2354837398	guest_token	BAhJIig1bmVjdW9ObnVzaHFpOE13SFBvY0lRMTcyMzY4NTM5ODY1NQY6BkVU--a8ea7a3c1feb345f279e3a0baccb74a02a835b5c
			#HttpOnly_catalog.ldc.upenn.edu	FALSE	/	TRUE	2354837398	token	BAhJIig1bmVjdW9ObnVzaHFpOE13SFBvY0lRMTcyMzY4NTM5ODY1NQY6BkVU--a8ea7a3c1feb345f279e3a0baccb74a02a835b5c
			#HttpOnly_.ldc.upenn.edu	TRUE	/	TRUE	1724895634	remember_spree_user_token	BAhbCFsGaQOYjgFJIhk1V0NGeVpSQ3pnZ1JlWXV5bTZmdQY6BkVUSSIWMTcyMzY4NjAzNC4xNDEyNDMGOwBG--e3c3ddcae8444aae4c86c79a420c8188177ddb8f
	*/
	// 验证Cookies并重定向下载
	view.OnDownload(func(url string) {
		if err := app.Downloader.DownloadFile(url, func(option *downloader.Option) {
			if app.Config.GetCookieFileABS() != "" {
				option.Cookies = app.Config.ParseCookie(app.Config.GetCookieFileABS())
			}
		}); err != nil {
			alert.Error("下载失败", err.Error())
			return
		}
		alert.Success("提示", "下载成功！")
	})

	app.KeepRunning()
}