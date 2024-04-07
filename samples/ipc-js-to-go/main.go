package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"time"

	blink "github.com/epkgs/mini-blink"
)

//go:embed static
var static embed.FS

func main() {
	app := blink.NewApp()
	defer app.Free()

	res, _ := fs.Sub(static, "static")
	app.Resource.Bind("local", res) // 将内嵌文件夹绑定到 FileSystem

	view := app.CreateWebWindowPopup(blink.WkeRect{
		W: 800, H: 800,
	})

	view.Window.MoveToCenter()

	view.LoadURL("http://local/index.html")

	view.OnDocumentReady(func(frame blink.WkeWebFrameHandle) {
		go timeTask(app)
	})

	view.ShowWindow()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	app.KeepRunning()
}

// 定时执行web js
func timeTask(app *blink.Blink) {
	//这里模拟go中触发js监听的事件
	var param0 = 0
	for {
		//每1秒钟执行一次
		time.Sleep(time.Second)
		fmt.Println("timeTask", param0)
		param0++
		//将数据发送出去
		app.IPC.Invoke("js-on-event-demo", fmt.Sprintf("Go发送的数据: %d", param0), float64(param0+10))
		// 如果JS返回结果, 需要通过回调函数入参方式接收返回值
		res := app.IPC.Invoke("js-on-event-demo-return", []interface{}{fmt.Sprintf("Go发送的数据: %d", param0), float64(param0 + 10)})

		//需要正确的获取类型，否则会失败
		fmt.Println("JS返回数据:", res.(string))
	}
}
