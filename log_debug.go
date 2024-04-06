//go:build !release

package blink

import (
	"fmt"
)

func _log(format string, vars ...any) {
	fmt.Printf(format, vars...)
	fmt.Print("\n")
}
