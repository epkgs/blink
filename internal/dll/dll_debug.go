//go:build !release

package dll

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

	dir := filepath.Join(root, "internal", "dll", "static")

	dirFS := os.DirFS(dir)

	return http.FS(dirFS)
}
