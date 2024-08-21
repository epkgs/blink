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

func pick(defaultTitle string, titleOrContent string, contents ...string) (title string, content string) {
	if len(contents) == 0 {
		return defaultTitle, titleOrContent
	} else {
		return titleOrContent, strings.Join(contents, "\r\n")
	}
}

func Alert(flag uint32, titleOrContent string, contents ...string) int32 {
	title, content := pick("注意", titleOrContent, contents...)
	return win.MessageBox(0, strToWcharPtr(content), strToWcharPtr(title), flag|win.MB_OK|win.MB_TOPMOST)
}

func Error(titleOrContent string, contents ...string) int32 {
	return Alert(win.MB_ICONERROR, titleOrContent, contents...)
}

func Info(titleOrContent string, contents ...string) int32 {
	return Alert(win.MB_ICONINFORMATION, titleOrContent, contents...)
}

func Warning(titleOrContent string, contents ...string) int32 {
	return Alert(win.MB_ICONWARNING, titleOrContent, contents...)
}
