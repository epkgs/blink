package alert

import (
	"strings"
	"syscall"

	"github.com/lxn/win"
)

func strToWcharPtr(s string) *uint16 {
	ptr, _ := syscall.UTF16PtrFromString(s)
	return ptr
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
	return Alert(win.MB_ICONERROR, titleOrContent, content...)
}

func Warn(titleOrContent string, content ...string) int32 {
	return Alert(win.MB_ICONWARNING, titleOrContent, content...)
}

func Success(titleOrContent string, content ...string) int32 {
	return Alert(win.MB_ICONINFORMATION, titleOrContent, content...)
}
