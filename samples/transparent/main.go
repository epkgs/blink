package main

import (
	"embed"
	"io/fs"
	"os"

	blink "github.com/epkgs/mini-blink"
)

//go:embed static
var static embed.FS

func main() {
	app := blink.NewApp()
	defer app.Free()

	res, _ := fs.Sub(static, "static")
	app.Resource.Bind("local", res) // 将内嵌文件夹绑定到 FileSystem

	view := app.CreateWebWindowTransparent(blink.WkeRect{
		W: 300, H: 300,
	})

	view.Window.MoveToCenter()

	view.LoadURL("http://local/transparent.html")

	view.ShowWindow()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	app.Run()
}
