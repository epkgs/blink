//go:build release && windows

package log

import (
	"fmt"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32")
var outputDebugStringW = kernel32.NewProc("OutputDebugStringW")

func log(format string, vars ...interface{}) {

	s := fmt.Sprintf(format, vars...)

	p, err := syscall.UTF16PtrFromString(s)
	if err == nil {
		outputDebugStringW.Call(uintptr(unsafe.Pointer(p)))
	}
}

func Debug(format string, vars ...interface{}) {
}
