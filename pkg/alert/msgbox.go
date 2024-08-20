package alert

import (
	"github.com/lxn/win"
	"strings"
	"syscall"
)

func strToWcharPtr(s string) *uint16 {
	ptr, _ := syscall.UTF16PtrFromString(s)
	return ptr
}

func messageBox(title, content string, flag uint32) int32 {
	hwnd := win.GetActiveWindow() // 当前激活的窗口
	win.SetForegroundWindow(hwnd) // 弹窗置于目标窗口之上
	return win.MessageBox(hwnd, strToWcharPtr(content), strToWcharPtr(title), flag)
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
	return messageBox(title, content, win.MB_OK|win.MB_ICONERROR)
}

func Info(titleOrContent string, contents ...string) int32 {
	title, content := pick("提示", titleOrContent, contents...)
	return messageBox(title, content, win.MB_OK|win.MB_ICONINFORMATION)
}

func Warning(titleOrContent string, contents ...string) int32 {
	title, content := pick("警告", titleOrContent, contents...)
	return messageBox(title, content, win.MB_OK|win.MB_ICONWARNING)
}
