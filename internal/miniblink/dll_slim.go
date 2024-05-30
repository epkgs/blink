//go:build slim

package miniblink

import (
	"errors"
	"io/fs"
)

const (
	ARCH    = ""
	VERSION = ""
)

type emptyFS struct {
}

func (fs emptyFS) Open(name string) (fs.File, error) {
	return nil, errors.New("slim 模式未嵌入 blink.dll")
}

var res = &emptyFS{}
