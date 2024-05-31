package util

import "sync"

type Pool[V any] struct {
	sync.Pool
	Reset func(data V)
}

func (pool *Pool[V]) GetAndReset() any {
	data := pool.Get().(V)
	resetFn := pool.Reset
	if resetFn != nil {
		resetFn(data)
	}
	return data
}
