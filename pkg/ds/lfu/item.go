package lfu

import (
	"homework/pkg/ds/linkedlist"
	"time"
)

type item[K comparable, V any] struct {
	expiredAt time.Time

	node *linkedlist.Node[K, V]
	freq int
}

func (i item[K, V]) Expired(now time.Time) bool {
	return i.expiredAt.Before(now)
}
