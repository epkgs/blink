package blink

func logInfo(format string, vars ...any) {
	_log("[GO-BLINK INFO] "+format, vars...)
}

func logWarning(format string, vars ...any) {
	_log("[GO-BLINK WARN] "+format, vars...)
}

func logError(format string, vars ...any) {
	_log("[GO-BLINK ERROR] "+format, vars...)
}
