//go:build !release

package blink

import (
	"github.com/epkgs/mini-blink/internal/devtools"
)

func (v *View) ShowDevTools() {

	if !Resource.IsExist("__devtools__") {
		Resource.Bind("__devtools__", devtools.FS)
	}

	var callback WkeOnShowDevtoolsCallback = func(hwnd WkeHandle, param uintptr) uintptr {

		view := NewView(v.mb, hwnd, WKE_WINDOW_TYPE_POPUP, v)

		view.ForceReload() // 必须刷新才会加载

		return 0
	}

	v.mb.CallFunc("wkeShowDevtools", uintptr(v.Hwnd), StringToWCharPtr("http://__devtools__/inspector.html"), CallbackToPtr(callback), 0)
}
