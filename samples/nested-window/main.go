package main

import (
	"embed"
	"io/fs"

	blink "github.com/epkgs/blink"
)

//go:embed web
var _web embed.FS
var webDir, _ = fs.Sub(_web, "web")

func main() {
	app := blink.NewApp()

	app.Resource.Bind("local", webDir)

	parent := app.CreateWebWindowPopup(blink.WithWebWindowSize(800, 600))
	parent.Window.EnableBorderResize(true)
	parent.Window.HideCaption()
	parent.Window.MoveToCenter()
	parent.OnDestroy(func() {
		app.Exit()
	})

	child := app.CreateWebWindowControl(parent,
		blink.WithWebWindowSize(800-4, 570-2),
		blink.WithWebWindowPos(2, 29),
	)

	parent.Window.OnSize(func(stype blink.SIZE_TYPE, width, height uint16) {
		child.Window.Resize(int32(width-4), int32(height-29-2))
	})

	parent.LoadURL("http://local/index.html")
	child.LoadURL("https://weixin.qq.com/")

	parent.ShowWindow()
	child.ShowWindow()

	app.KeepRunning()
}
