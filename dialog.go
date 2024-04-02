package blink

import (
	"strings"

	"github.com/lxn/win"
)

func MessageBox(hwnd uintptr, title, content string, flags ...uint32) int32 {

	var flag uint32 = 0
	for _, f := range flags {
		flag = flag | f
	}

	return win.MessageBox(win.HWND(hwnd), StringToWcharU16Ptr(content), StringToWcharU16Ptr(title), flag)
}

func msgbox(hwnd uintptr, flags uint32, titleOrContent string, lines ...string) int32 {
	var title string
	var content string
	if len(lines) == 0 {
		title = "错误"
		content = titleOrContent
	} else {
		title = titleOrContent
		content = strings.Join(lines, "\n")
	}

	return MessageBox(hwnd, title, content, flags)
}

func MessageBoxError(hwnd uintptr, titleOrContent string, lines ...string) int32 {
	return msgbox(hwnd, win.MB_OK|win.MB_ICONERROR, titleOrContent, lines...)
}
