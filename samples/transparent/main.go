package main

import (
	"os"

	blink "github.com/epkgs/mini-blink"
)

func main() {
	app := blink.NewApp()
	defer app.Free()

	blink.Resource.Bind("local", "static") // 将本地文件绑定到 FileSystem

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
