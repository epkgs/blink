//go:build !slim && 386

package dll

import (
	"embed"
	"io/fs"
)

//go:embed x32
var res embed.FS

const DLL_FILE = "blink.dll"

var FS, _ = fs.Sub(res, "x32")
