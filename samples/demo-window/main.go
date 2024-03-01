package main

import (
	"fmt"
	"os"
	"path/filepath"

	blink "github.com/epkgs/mini-blink"
)

func main() {
	app := blink.NewApp()
	defer app.Free()

	pwd, _ := os.Getwd()
	dir := filepath.Join(pwd, "cmd", "demo-window", "static") // ! 默认是从项目根目录开始检索，由于demo目录不是项目根目录，所以需要配置绝对路径
	blink.Resource.BindDir("local", dir)                      // 将本地文件绑定到 FileSystem

	view := app.CreateWebWindowTransparent(blink.WkeRect{
		W: 800, H: 800,
	})

	view.Window.MoveToCenter()

	view.EnableBorderResize()

	view.LoadURL("http://local/window.html")

	view.Show()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	view.AddEventListener(".custom-zone", "mouseover", func() {
		fmt.Printf("custom zone hover\n")
	})

	app.KeepRunning()
}
