package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	blink "github.com/epkgs/blink"
)

//go:embed static
var static embed.FS

func main() {
	app := blink.NewApp()
	defer app.Exit()

	res, _ := fs.Sub(static, "static")
	app.Resource.Bind("local", res) // 将内嵌文件夹绑定到 FileSystem

	//一个普通的窗体
	view := app.CreateWebWindowPopup(func(c *blink.WebWindowConfig) {
		c.W = 800
		c.H = 600
	})

	view.Window.SetTitle("JS互操作")
	view.Window.SetIconFromFile("app.ico") // 相对路径默认从项目根目录开始
	view.Window.MoveToCenter()

	view.LoadURL("http://local/call_js.html")
	view.ShowWindow()

	view.OnConsole(func(level int, message, sourceName string, sourceLine int, stackTrace string) {
		fmt.Printf("js console: %s\n", message)
	})

	view.OnDocumentReady(func(frame blink.WkeWebFrameHandle) {
		// 避免阻塞主线程
		go func() {
			// 调用func_1
			view.CallJsFunc("func_1", "张三", 18)

			// 等待返回值
			resp2, err := view.CallJsFunc("func_2").Wait()
			if err != nil {
				fmt.Printf("call js func_2 error: %s\n", err)
				return
			}
			result2 := (resp2).(string) // 断言为字符串
			fmt.Printf("func_2 result is %s\n", result2)

			//获取func_3返回的非基本数据类型
			resp3, err := view.CallJsFunc("func_3").Wait()
			if err != nil {
				fmt.Printf("call js func_3 error: %s\n", err)
				return
			}
			result3 := (resp3).(map[string]interface{}) // wait for result
			fmt.Printf("func_3 result is %v\n", result3)
		}()
	})

	view.OnDestroy(func() {
		os.Exit(0)
	})

	view.ShowDevTools()

	app.KeepRunning()
}
