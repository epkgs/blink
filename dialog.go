package blink

import "github.com/lxn/win"

func MessageBox(hwnd uintptr, text, caption string, flags uint32) int32 {
	return win.MessageBox(win.HWND(hwnd), StringToWcharU16Ptr(text), StringToWcharU16Ptr(caption), flags)
}

func MessageBoxError(hwnd uintptr, text, caption string) int32 {
	return MessageBox(hwnd, text, caption, win.MB_OK|win.MB_ICONERROR)
}
