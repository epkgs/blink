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

	view := app.CreateWebWindowTransparent(func(c *blink.WebWindowConfig) {
		c.W = 800
		c.H = 600
	})

	view.Window.MoveToCenter()

	view.LoadURL("http://local/flameless.html")

	view.ShowWindow()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	view.OnDocumentReady(func(frame blink.WkeWebFrameHandle) {
		view.AddEventListener(".custom-zone", "mouseover", func() {
			fmt.Printf("custom zone hover\n")
		})
	})

	app.KeepRunning()
}
