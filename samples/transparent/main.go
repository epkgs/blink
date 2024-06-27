package main

import (
	"embed"
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

	view := app.CreateWebWindowTransparent(func(config *blink.WebWindowConfig) {
		config.W = 300
		config.H = 300
	})

	view.Window.EnableBorderResize(false)
	view.Window.MoveToCenter()

	view.LoadURL("http://local/transparent.html")

	view.ShowWindow()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	app.KeepRunning()
}
