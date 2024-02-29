//go:generate goversioninfo
package main

import (
	"os"

	blink "github.com/epkgs/mini-blink"
)

func main() {
	app := blink.NewApp()
	defer app.Free()

	view := app.CreateWebWindowPopup()
	view.Window.SetIconFromBytes(icon)
	view.Window.SetTitle("miniblink窗口")
	view.Window.MoveToCenter()
	view.LoadURL("https://www.baidu.com")
	view.Show()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	app.KeepRunning()
}
