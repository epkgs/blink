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

		app.AddJob(func() {

			fmt.Println("timeTask", param0)
			param0++
			//将数据发送出去
			app.IPC.Invoke("js-on-event-demo", param0, float64(param0+10), "this is a test")
			res, _ := app.IPC.Invoke("js-on-event-demo-return", param0, float64(param0+10))

			fmt.Printf("JS返回数据: %v\n", res) // ! 如需要正确的获取类型，请注意断言正确类型，否则将会导致 panic
		})
	}
}
