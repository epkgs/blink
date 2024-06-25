package alert

import (
	"strings"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

func strToWcharPtr(s string) *uint16 {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		*p = 0
	}
	return p
}

func messageBox(hwnd uintptr, title, content string, flag uint32) int32 {
	return win.MessageBox(win.HWND(hwnd), strToWcharPtr(content), strToWcharPtr(title), flag)
}

func Alert(flags uint32, titleOrContent string, content ...string) int32 {
	var title string
	var txt string
	if len(content) == 0 {
		title = "错误"
		txt = titleOrContent
	} else {
		title = titleOrContent
		txt = strings.Join(content, "\n")
	}

	return messageBox(0, title, txt, flags)
}

func Error(titleOrContent string, content ...string) int32 {
	return Alert(win.MB_OK|win.MB_ICONERROR, titleOrContent, content...)
}
