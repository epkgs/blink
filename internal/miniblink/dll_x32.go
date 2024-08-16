//go:build !slim && 386

package miniblink

import (
	"embed"
)

const (
	ARCH    = "x32"
	VERSION = "4975"
	DllHash = "073889b81f42989cd4bb71613e03b3c5214f4bab5185c6702b6fc647725db0ec"
)

//go:embed release/x32
var res embed.FS

/*
windows sha256
certutil -hashfile .\release\x32\miniblink_4975_x32.dll sha256 | Select-Object -Skip 1 -First 1
*/
