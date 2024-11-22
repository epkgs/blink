package utils

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func Go(f func(), handleError func(error)) {
	go func() {

		defer func() {
			if r := recover(); r != nil {

				if err, ok := r.(error); ok {
					if handleError != nil {
						handleError(err)
						return
					}
				}

				log.Printf("panic by goroutine: %v", r)
			}
		}()

		f()

	}()
}

func GoWithContext(ctx context.Context, f func(), handleError func(error)) {
	go func() {

		defer func() {
			if r := recover(); r != nil {

				if err, ok := r.(error); ok {
					if handleError != nil {
						handleError(err)
						return
					}
				}

				log.Printf("panic by goroutine: %v\n", r)
			}
		}()

		var once sync.Once
		done := false

		// 不断循环，检查 ctx 是否已经结束
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// 保证仅执行一次
				once.Do(func() {
					f()
					done = true
				})

				// 判断是否已经执行完毕，执行完毕跳出整个协程
				if done {
					return
				}
			}
		}
	}()
}
