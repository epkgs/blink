package log

import "fmt"

func log(format string, vars ...interface{}) {
	fmt.Printf(format+"\n", vars...)
}

func Debug(format string, vars ...interface{}) {
	log("[GO-BLINK DEBUG] "+format, vars...)
}

func Info(format string, vars ...interface{}) {
	log("[GO-BLINK INFO] "+format, vars...)
}

func Warning(format string, vars ...interface{}) {
	log("[GO-BLINK WARN] "+format, vars...)
}

func Error(format string, vars ...interface{}) {
	log("[GO-BLINK ERROR] "+format, vars...)
}
