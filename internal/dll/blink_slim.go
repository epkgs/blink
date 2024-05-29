//go:build slim

package dll

import (
	"errors"
	"io/fs"
)

type emptyFS struct {
}

func (fs emptyFS) Open(name string) (fs.File, error) {
	return nil, errors.New("slim 模式未嵌入 blink.dll")
}

var _fs = &emptyFS{}

var FS = fs.FS(_fs)
