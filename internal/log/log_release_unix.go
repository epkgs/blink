//go:build release && !windows

package log

func log(format string, vars ...interface{}) {
}

func Debug(format string, vars ...interface{}) {
}
