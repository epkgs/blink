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
	view.ShowWindow()

	pwd, _ := os.Getwd()

	var stop func()
	stop = view.OnDocumentReady(func(frame blink.WkeWebFrameHandle) {
		go func() {
			if !view.IsMainFrame(frame) {
				return
			}
			if stop != nil {
				stop()
			}
			// 等待图片加载完成
			time.Sleep(time.Second * 3)
			// 保存为pdf
			view.SaveToPDF(path.Join(pwd, "screenshot.pdf"))
			os.Exit(0)
		}()
	})

	view.OnDestroy(func() {
		os.Exit(0)
	})

	app.KeepRunning()
}
