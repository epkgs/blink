//go:build release && !windows

package blink

func _log(format string, vars ...any) {
}
