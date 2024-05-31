package collection

import (
	"time"
)

type ExpiredHashMap[K comparable, V any] struct {
	data map[K]V
}

func NewExpiredHashMap[K comparable, V any]() *ExpiredHashMap[K, V] {
	return &ExpiredHashMap[K, V]{data: make(map[K]V, 16)}
}

func (hashMap *ExpiredHashMap[K, V]) Put(key K, value V, expiredTime time.Duration) {
	hashMap.data[key] = value
	time.AfterFunc(expiredTime, func() {
		delete(hashMap.data, key)
	})
}

func (hashMap *ExpiredHashMap[K, V]) Get(key K) (value V, ok bool) {
	value, ok = hashMap.data[key]
	return
}
