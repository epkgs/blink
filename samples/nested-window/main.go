package main

import (
	"embed"
	"io/fs"

	blink "github.com/epkgs/mini-blink"
	"github.com/lxn/win"
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

	parent.Window.OnResize(func(r *win.RECT) {
		width := r.Right - r.Left - 4
		height := r.Bottom - r.Top - 29 - 2

		// child 的 x, y 坐标是相对于 parent 的
		child.Window.Resize(2, 29, width, height)
	})

	parent.LoadURL("http://local/index.html")
	child.LoadURL("https://weixin.qq.com/")

	parent.ShowWindow()
	child.ShowWindow()

	app.KeepRunning()
}
