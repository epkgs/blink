package blink

import (
	"fmt"
)

func _log(format string, vars ...any) {
	fmt.Printf(format, vars...)
}

func logInfo(format string, vars ...any) {
	_log("[INFO] "+format, vars...)
}

func logWarning(format string, vars ...any) {
	_log("[WARN] "+format, vars...)
}

func logError(format string, vars ...any) {
	_log("[ERROR] "+format, vars...)
}
