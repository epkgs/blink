//go:build !release

package blink

import (
	"unsafe"
)

var env = &Env{
	isSYS64: unsafe.Sizeof(uintptr(0)) == 8,
	isDebug: true,
}
