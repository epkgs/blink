//go:build !release

package log

import (
	"fmt"
)

func log(format string, vars ...any) {
	fmt.Printf(format, vars...)
	fmt.Print("\n")
}
