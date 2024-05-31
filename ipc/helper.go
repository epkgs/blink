package ipc

import blink "github.com/epkgs/mini-blink"

func Sent(mb *blink.Blink, channel string, args ...interface{}) error {
	return mb.IPC.Sent(channel, args...)
}

func Invoke[R interface{}](mb *blink.Blink, channel string, args ...interface{}) (R, error) {
	res, err := mb.IPC.Invoke(channel, args...)
	return res.(R), err
}

func Handle(mb *blink.Blink, channel string, handler blink.Callback) error {
	return mb.IPC.Handle(channel, handler)
}
