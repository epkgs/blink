package queue

import (
	"sync"
)

type IQueue[T interface{}] interface {
	First() T
	Last() T
}

type Queue[T interface{}] struct {
	ch   chan IQueue[T]
	list []T
	mu   sync.Mutex
}

func NewQueue[T interface{}](length ...int) *Queue[T] {
	count := 100

	if len(length) > 0 {
		count = length[0]
	}

	return &Queue[T]{
		ch:   make(chan IQueue[T], count),
		list: []T{},
	}
}

// 队列末尾添加新的元素
func (q *Queue[T]) AddLast(v T) *Queue[T] {

	q.mu.Lock()
	defer q.mu.Unlock()

	q.list = append(q.list, v)

	q.ch <- q

	return q
}

// 队列开头添加新的元素
func (q *Queue[T]) AddFirst(v T) *Queue[T] {

	q.mu.Lock()
	defer q.mu.Unlock()

	q.list = append([]T{v}, q.list...)
	q.ch <- q

	return q
}

func (q *Queue[T]) Chan() chan IQueue[T] {

	return q.ch
}

// 删除第一个元素并返回
func (q *Queue[T]) First() T {

	q.mu.Lock()
	defer q.mu.Unlock()

	first := q.list[0]
	q.list = append(q.list[:0], q.list[1:]...)

	return first
}

// 删除最后一个元素并返回
func (q *Queue[T]) Last() T {

	q.mu.Lock()
	defer q.mu.Unlock()

	lastIdx := len(q.list) - 1
	last := q.list[lastIdx]
	q.list = q.list[:lastIdx]

	return last
}
