//go:build !386

package dll

import (
	"embed"
	"io/fs"
)

//go:embed x64
var res embed.FS

var FS, _ = fs.Sub(res, "x64")
