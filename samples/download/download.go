package main

import (
	"compress/flate"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/epkgs/blink"
	"github.com/epkgs/blink/pkg/downloader"
	"github.com/gin-gonic/gin"
)

func main() {
	app := blink.NewApp()

	go runWebServer()

	// 单线程直接下载
	_, _ = app.Download("https://httpbin.org/robots.txt")

	// deflate压缩内容下载
	_, _ = app.Download("https://comment.bilibili.com/177987845.xml", func(c *downloader.Config) {
		c.Interceptors.HttpDownloading = func(job *downloader.Job, res *http.Response) io.Reader {
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
	_, _ = app.Download("https://httpbin.org/image", func(c *downloader.Config) {
		c.Interceptors.BeforeSaveFile = func(job *downloader.Job) {
			// 修改文件名
			job.FileName += ".png"
		}
	})

	// 无法 HEAD 的文件下载
	_, _ = app.Download("http://localhost:9999/d1/test.zip")

	// 登录后才能下载
	view := app.CreateWebWindowPopup()
	view.LoadURL("http://localhost:9999/login/loginedUser") // 使用网页加载，使其能有cookie
	view.OnDocumentReady(func(frame blink.WkeWebFrameHandle) {
		defer app.Exit()
		_, _ = app.Download("http://localhost:9999/d2/test.zip")
	})

	app.KeepRunning()
}

func runWebServer() {
	// 使用gin创建一个后端下载接口
	router := gin.Default()

	pwd, _ := os.Getwd()
	staticPath := path.Join(pwd, "static")

	router.Static("/static", staticPath)

	// 设置登录路由
	router.GET("/login/:auth", loginHandler)

	router.GET("/d1/:filename", func(c *gin.Context) { // head为404
		filename := c.Param("filename")
		c.File(path.Join(staticPath, filename))
	})

	router.GET("/d2/:filename", func(c *gin.Context) { // head为404
		// 检查 cookies 中是否存在 "auth" cookie
		cookie, err := c.Cookie("auth")
		if err != nil || cookie != "true" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}
		filename := c.Param("filename")
		changedName := "mytest.zip"
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; fileName=%s", changedName))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(path.Join(staticPath, filename))
	})

	// 启动服务器监听在指定端口
	router.Run(":9999")

}

// 登录处理函数
func loginHandler(c *gin.Context) {
	auth := c.Param("auth")
	if auth != "" {
		// 登录成功，设置 cookie
		c.SetCookie("auth", "true", 3600, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "Login successful"})
}
