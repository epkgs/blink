package main

import (
	"encoding/json"
	"fmt"
	"os"

	blink "github.com/epkgs/mini-blink"
)

func main() {
	app := blink.NewApp()
	defer app.Free()

	blink.Resource.Bind("local", "static") // 将本地文件绑定到 FileSystem

	//一个普通的窗体
	view := app.CreateWebWindowPopup(blink.WkeRect{
		W: 800,
		H: 500,
	})

	view.Window.SetTitle("JS互操作")
	view.Window.SetIconFromFile("app.ico") // 相对路径默认从项目根目录开始
	view.Window.MoveToCenter()

	view.LoadURL("http://local/call_js.html")
	view.Show()

	view.OnConsole(func(level int, message, sourceName string, sourceLine int, stackTrace string) {
		fmt.Printf("js console: %s\n", message)
	})

	view.OnDocumentReady(func(frame blink.WkeWebFrameHandle) {
		//调用func_1
		view.CallJsFunc(nil, "func_1", "张三", 18)

		//获取func_2返回的基础数据类型
		view.CallJsFunc(func(result any) {

			fmt.Printf("func_2 result is %s\n", result.(string))
		}, "func_2")

		//获取func_3返回的非基本数据类型
		view.CallJsFunc(func(result any) {
			bytes, _ := json.Marshal(result.(map[string]any))
			fmt.Printf("func_3 result is %s\n", string(bytes))
		}, "func_3")
	})

	view.OnDestroy(func() {
		os.Exit(0)
	})

	view.ShowDevTools()

	app.KeepRunning()
}
