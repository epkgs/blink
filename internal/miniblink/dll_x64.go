//go:build !slim && !386

package miniblink

import (
	"embed"
)

const (
	ARCH    = "x64"
	VERSION = "4975"
	DllHash = "2e5a7b64260461012a195eaf6e1dc417b011932b6fc232849653f4b1693564d1"
)

//go:embed release/x64
var res embed.FS

/*
windows sha256
certutil -hashfile .\release\x64\miniblink_4975_x64.dll sha256 | Select-Object -Skip 1 -First 1
*/
