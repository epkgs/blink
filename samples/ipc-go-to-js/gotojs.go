package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"

	blink "github.com/epkgs/blink"
)

//go:embed resources
var resources embed.FS

func main() {

	app := blink.NewApp()
	defer app.Exit()

	res, _ := fs.Sub(resources, "resources")
	app.Resource.Bind("local", res) // 将内嵌文件夹绑定到 FileSystem

	view := app.CreateWebWindowPopup(func(c *blink.WebWindowConfig) {
		c.W = 800
		c.H = 600
	})

	view.Window.MoveToCenter()

	view.LoadURL("http://local/go-to-js.html")

	view.ShowWindow()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	//在go中监听一个事件, 不带返回值
	//使用上下文获取参数
	app.IPC.Handle("go-on-event-demo", func(arg1 string, arg2 string, arg3 int) {
		//参数是以js调用时传递的参数下标位置开始计算，从0开始表示第1个参数
		fmt.Printf("[go-on-event-demo] arg1: '%s', arg2: '%s', arg3: %d\n", arg1, arg2, arg3)
	})

	//带有返回值的事件
	app.IPC.Handle("go-on-event-demo-return", func(p1 string, p2 int, p3 bool, p4 float64, p5 string) string {
		fmt.Println("go-on-event-demo-return event run")
		fmt.Println("\t参数1-length:", len(p1), p1)
		//fmt.Println("\t参数1:", p1)
		fmt.Println("\t参数2:", p2)
		fmt.Println("\t参数3:", p3)
		fmt.Println("\t参数4:", p4)
		fmt.Println("\t参数5:", p5)
		//返回给JS数据, 通过 context.Result()
		var buf = bytes.Buffer{}
		for i := 100; i < 200; i++ {
			buf.WriteString(fmt.Sprintf("[%d]-", i))
		}
		var data = "这是在GO中监听事件返回给JS的数据:" + buf.String()
		fmt.Println("返回给JS数据 - length:", strings.Count(data, "")-1)
		return data
	})

	// 在Go中监听一个事件, 不带返回值
	// 使用形参接收参数
	// 在JS中入参类型必须相同
	app.IPC.Handle("go-on-event-demo-argument", func(p1 int, p2 string, p3 float64, p4 bool, p5 string) {
		fmt.Println("param1:", p1)
		fmt.Println("param2:", p2)
		fmt.Println("param3:", p3)
		fmt.Println("param4:", p4)
		fmt.Println("param5:", p5)
	})

	// 在Go中监听一个事件, 带返回值
	// 使用形参接收参数
	// 在JS中入参类型必须相同
	// 返回参数可以同时返回多个, 在JS接收时同样使用回调函数方式以多个入参形式接收
	app.IPC.Handle("go-on-event-demo-argument-return", func(p1 int, p2 string, p3 float64, p4 bool, p5 string) string {
		fmt.Println("param1:", p1)
		fmt.Println("param2:", p2)
		fmt.Println("param3:", p3)
		fmt.Println("param4:", p4)
		fmt.Println("param5:", p5)
		return fmt.Sprintf("%d-%s-%f-%t-%s", p1, p2, p3, p4, p5)
	})

	app.KeepRunning()
}
