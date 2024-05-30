//go:build release && !386

package blink

var env = &Env{
	isSYS64: true,
	isDebug: false,
}
