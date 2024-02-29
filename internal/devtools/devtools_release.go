//go:build release

package devtools

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed front_end
var frontEnd embed.FS

func AssetFile() http.FileSystem {
	devtools, err := fs.Sub(frontEnd, "front_end")

	if err != nil {
		panic(err)
	}
	return http.FS(devtools)
}
