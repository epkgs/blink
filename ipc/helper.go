package ipc

import "github.com/epkgs/mini-blink"

func Sent(mb *blink.Blink, channel string, args ...any) error {
	return mb.IPC.Sent(channel, args...)
}

func Invoke[R any](mb *blink.Blink, channel string, args ...any) (R, error) {
	res, err := mb.IPC.Invoke(channel, args...)
	return res.(R), err
}

func Handle(mb *blink.Blink, channel string, handler blink.Callback) error {
	return mb.IPC.Handle(channel, handler)
}
