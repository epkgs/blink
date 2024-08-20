//go:build !slim && 386

package miniblink

import (
	"embed"
)

const (
	ARCH    = "x32"
	VERSION = "4975"
	// windows产生sha256命令：certutil -hashfile .\internal\miniblink\release\x32\miniblink_4975_x32.dll sha256 | Select-Object -Skip 1 -First 1
	DllSHA256 = "08eeb6e4f5e80eff17d72f5abbf8fe14f66777db81fa4fa3101c16c3c3df3409"
)

//go:embed release/x32
var res embed.FS
