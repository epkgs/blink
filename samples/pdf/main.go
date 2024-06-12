//go:generate goversioninfo
package main

import (
	"os"
	"path"
	"time"

	blink "github.com/epkgs/mini-blink"
)

func main() {
	app := blink.NewApp()
	defer app.Free()

	view := app.CreateWebWindowPopup()
	view.Window.SetIconFromBytes(icon)
	view.Window.SetTitle("miniblink窗口")
	view.Window.MoveToCenter()
	view.LoadURL("https://www.baidu.com")
	view.ShowWindow() // 可不显示窗口，静默进行

	pwd, _ := os.Getwd()

	// 初始化 stop 函数为空函数
	var stop func() = func() {}
	stop = view.OnDocumentReady(func(frame blink.WkeWebFrameHandle) {
		// 在协程内处理，避免 sleep 阻塞主线程
		go func() {
			// 仅处理主frame，忽略其他frame的 document ready 事件
			if !view.IsMainFrame(frame) {
				return
			}
			// 已触发 document ready，取消监听
			stop()
			// 等待图片加载完成
			time.Sleep(time.Second * 3)
			// 保存为pdf
			view.SaveWebFrameToPDF(frame, path.Join(pwd, "screenshot.pdf"))
			// 截图完成，退出
			os.Exit(0)
		}()
	})

	view.OnDestroy(func() {
		os.Exit(0)
	})

	app.KeepRunning()
}
