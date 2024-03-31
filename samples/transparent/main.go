package main

import (
	"os"
	"path/filepath"

	blink "github.com/epkgs/mini-blink"
)

func main() {
	app := blink.NewApp()
	defer app.Free()

	pwd, _ := os.Getwd()
	dir := filepath.Join(pwd, "samples", "transparent", "static") // ! 默认是从项目根目录开始检索，由于demo目录不是项目根目录，所以需要配置绝对路径
	blink.Resource.BindDir("local", dir)                          // 将本地文件绑定到 FileSystem

	view := app.CreateWebWindowTransparent(blink.WkeRect{
		W: 300, H: 300,
	})

	view.Window.MoveToCenter()

	view.LoadURL("http://local/transparent.html")

	view.Show()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	app.KeepRunning()
}
