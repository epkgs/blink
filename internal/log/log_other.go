//go:build !release || !windows

package log

import "fmt"

func log(format string, vars ...interface{}) {
	fmt.Printf(format+"\n", vars...)
}
