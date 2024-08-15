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

func pick(defaultTitle string, titleOrContent string, contents ...string) (title string, content string) {
	if len(contents) == 0 {
		return defaultTitle, titleOrContent
	} else {
		return titleOrContent, strings.Join(contents, "\r\n")
	}
}

func Error(titleOrContent string, contents ...string) int32 {
	title, content := pick("错误", titleOrContent, contents...)
	return messageBox(0, title, content, win.MB_OK|win.MB_ICONERROR)
}

func Info(titleOrContent string, contents ...string) int32 {
	title, content := pick("提示", titleOrContent, contents...)
	return messageBox(0, title, content, win.MB_OK|win.MB_ICONINFORMATION)
}

func Warning(titleOrContent string, contents ...string) int32 {
	title, content := pick("警告", titleOrContent, contents...)
	return messageBox(0, title, content, win.MB_OK|win.MB_ICONWARNING)
}
