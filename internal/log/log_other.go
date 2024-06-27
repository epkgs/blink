//go:build !release && !debug

package log

import (
	"fmt"
)

func log(format string, vars ...interface{}) {
	fmt.Printf(format+"\n", vars...)
}

func Debug(format string, vars ...interface{}) {
	log("[GO-BLINK DEBUG] "+format, vars...)
}
