//go:build !release

package devtools

import (
	"net/http"
	"os"
	"path/filepath"
)

func AssetFile() http.FileSystem {

	root, err := os.Getwd()
	if err != nil {
		panic("获取运行目录出错")
	}

	dir := filepath.Join(root, "internal", "devtools", "front_end")

	dirFS := os.DirFS(dir)

	return http.FS(dirFS)
}
