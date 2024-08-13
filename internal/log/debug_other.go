//go:build !release || (release && debug)

package log

func Debug(format string, vars ...interface{}) {
	log("[GO-BLINK DEBUG] "+format, vars...)
}
