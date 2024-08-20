//go:build !slim && !386

package miniblink

import (
	"embed"
)

const (
	ARCH    = "x64"
	VERSION = "4975"
	// windows产生sha256命令：certutil -hashfile .\internal\miniblink\release\x64\miniblink_4975_x64.dll sha256 | Select-Object -Skip 1 -First 1
	DllSHA256 = "82e32017b9d5832ff0f9c807bf219d4638d5a174da8637720e209971796dac95"
)

//go:embed release/x64
var res embed.FS
