package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	blink "github.com/epkgs/mini-blink"
)

//go:embed resources
var resources embed.FS

func main() {

	app := blink.NewApp()
	defer app.Free()

	res, _ := fs.Sub(resources, "resources")
	app.Resource.Bind("local", res) // 将内嵌文件夹绑定到 FileSystem

	view := app.CreateWebWindowPopup(blink.WkeRect{
		W: 800, H: 800,
	})

	view.Window.MoveToCenter()

	view.LoadURL("http://local/go-to-js.html")

	view.ShowWindow()

	//在go中监听一个事件, 不带返回值
	//使用上下文获取参数
	app.IPC.Handle("go-on-event-demo", func(args ...any) any {
		fmt.Println("go-on-event-demo event run")
		//js 中传递的数据
		fmt.Println("参数个数:", len(args))
		//参数是以js调用时传递的参数下标位置开始计算，从0开始表示第1个参数
		p1 := args[0]
		fmt.Println("参数1:", p1)

		return nil
	})

	//带有返回值的事件
	app.IPC.Handle("go-on-event-demo-return", func(args ...any) any {
		fmt.Println("go-on-event-demo-return event run")
		//js 中传递的数据
		fmt.Println("参数个数:", len(args))
		//参数是以js调用时传递的参数下标位置开始计算，从0开始表示第1个参数
		p1 := args[0].(string)
		p2 := int(args[1].(float64))
		p3 := args[2].(bool)
		p4 := args[3].(float64)
		p5 := args[4].(string)
		fmt.Println("\t参数1-length:", len(p1), p1)
		//fmt.Println("\t参数1:", p1)
		fmt.Println("\t参数2:", p2)
		fmt.Println("\t参数3:", p3)
		fmt.Println("\t参数4:", p4)
		fmt.Println("\t参数5:", p5)
		//返回给JS数据, 通过 context.Result()
		var buf = bytes.Buffer{}
		for i := 0; i < 100000; i++ {
			buf.WriteString(fmt.Sprintf("[%d]-", i))
		}
		var data = "这是在GO中监听事件返回给JS的数据:" + buf.String()
		fmt.Println("返回给JS数据 - length:", strings.Count(data, "")-1)
		return data
	})

	// 在Go中监听一个事件, 不带返回值
	// 使用形参接收参数
	// 在JS中入参类型必须相同
	app.IPC.Handle("go-on-event-demo-argument", func(args ...any) any {
		fmt.Println("param1:", args[0])
		fmt.Println("param2:", args[1])
		fmt.Println("param3:", int(args[2].(float64)))
		fmt.Println("param4:", args[3])
		fmt.Println("param5:", args[4])

		return nil
	})

	// 在Go中监听一个事件, 带返回值
	// 使用形参接收参数
	// 在JS中入参类型必须相同
	// 返回参数可以同时返回多个, 在JS接收时同样使用回调函数方式以多个入参形式接收
	app.IPC.Handle("go-on-event-demo-argument-return", func(args ...any) any {
		fmt.Println("param1:", args[0])
		fmt.Println("param2:", args[1])
		fmt.Println("param3:", int(args[2].(float64)))
		fmt.Println("param4:", args[3])
		fmt.Println("param5:", args[4])
		return fmt.Sprintf("%f-%v-%v-%v-%v", args[0], args[1], args[2], args[3], args[4])
	})

	app.KeepRunning()
}
