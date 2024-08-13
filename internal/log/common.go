package log

func Info(format string, vars ...interface{}) {
	log("[GO-BLINK INFO] "+format, vars...)
}

func Warning(format string, vars ...interface{}) {
	log("[GO-BLINK WARN] "+format, vars...)
}

func Error(format string, vars ...interface{}) {
	log("[GO-BLINK ERROR] "+format, vars...)
}
