package blink

import "unsafe"

type Env struct {
	isSYS64 bool
	isDebug bool
}

var env = &Env{
	isSYS64: unsafe.Sizeof(uintptr(0)) == 8,
	isDebug: true,
}

func IsSYS64() bool {
	return env.isSYS64
}

func IsDebug() bool {
	return env.isDebug
}

func IsRelease() bool {
	return !env.isDebug
}
