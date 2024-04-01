package main

import (
	"fmt"
	"os"

	blink "github.com/epkgs/mini-blink"
)

func main() {
	app := blink.NewApp()
	defer app.Free()

	blink.Resource.Bind("local", "static") // 将本地文件绑定到 FileSystem

	view := app.CreateWebWindowTransparent(blink.WkeRect{
		W: 800, H: 800,
	})

	view.Window.MoveToCenter()

	view.EnableBorderResize()

	view.LoadURL("http://local/no-border.html")

	view.Show()

	view.OnDestroy(func() {
		os.Exit(0)
	})

	view.AddEventListener(".custom-zone", "mouseover", func() {
		fmt.Printf("custom zone hover\n")
	})

	app.KeepRunning()
}
