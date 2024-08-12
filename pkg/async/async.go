package async

import (
	"fmt"
	"time"
)

type task[T any, F func() (T, error)] struct {
	Timeout   time.Duration
	fn        F
	isStarted bool
	resultCh  chan T
	errChan   chan error
}

type Pendding[T any] interface {
	Start() InProgress[T]
	Wait() (T, error)
}

type InProgress[T any] interface {
	Wait() (T, error)
}

func New[T any, F func() (T, error)](timeout time.Duration, fn F) Pendding[T] {
	return &task[T, F]{
		Timeout:  timeout,
		fn:       fn,
		resultCh: make(chan T, 1),
		errChan:  make(chan error, 1),
	}
}

func (f *task[T, F]) Start() InProgress[T] {

	if f.isStarted {
		return f
	}
	f.isStarted = true

	go func() {

		select {
		case <-time.After(f.Timeout):
			f.errChan <- fmt.Errorf("timeout after %v", f.Timeout.String())
			return
		default:
			res, err := f.fn()
			if err != nil {
				f.errChan <- err
			} else {
				f.resultCh <- res
			}
		}

	}()
	return f
}

func (f *task[T, F]) Wait() (T, error) {

	// 如果未启动，则启动
	if !f.isStarted {
		f.Start()
	}

	select {
	case err := <-f.errChan:
		return *new(T), err
	case res := <-f.resultCh:
		return res, nil
	}
}
