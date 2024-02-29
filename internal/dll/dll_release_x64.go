//go:build release && !386

package dll

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/blink_x64.dll
var static embed.FS

func AssetFile() http.FileSystem {
	res, err := fs.Sub(static, "static")

	if err != nil {
		panic(err)
	}
	return http.FS(res)
}
