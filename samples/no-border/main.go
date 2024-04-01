package main

import (
	"embed"
	"fmt"
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
	blink.Resource.Bind("local", res) // 将内嵌文件夹绑定到 FileSystem

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
