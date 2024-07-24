package lru

import (
	"homework/pkg/ds/linkedlist"
	"time"
)

type Cache[K comparable, V any] struct {
	ttl      time.Duration
	data     map[K]*item[K, V]
	list     *linkedlist.List[K, V]
	capacity int

	defaultV V
}

func NewLRUCache[K comparable, V any](capacity int) Cache[K, V] {
	return Cache[K, V]{
		list:     linkedlist.New[K, V](),
		data:     make(map[K]*item[K, V]),
		capacity: capacity,
	}
}

func (LRU *Cache[K, V]) Get(key K) (V, bool) {
	node, ok := LRU.data[key]
	if !ok {
		var res V
		return res, false
	}

	LRU.list.DeleteNode(node.node)
	if node.Expired(time.Now()) {
		delete(LRU.data, node.node.GetKey())
		return LRU.defaultV, false
	}

	LRU.data[key] = &item[K, V]{
		node:      LRU.list.PushNode(linkedlist.NewNode[K, V](key, node.node.GetValue())),
		expiredAt: time.Now().Add(LRU.ttl),
	}

	return node.node.GetValue(), true
}

func (LRU *Cache[K, V]) Put(key K, value V) {
	if node, ok := LRU.data[key]; ok {
		LRU.list.DeleteNode(node.node)
	}

	node := linkedlist.NewNode[K, V](key, value)
	LRU.list.PushNode(node)

	item := &item[K, V]{
		expiredAt: time.Now().Add(LRU.ttl),
		node:      node,
	}

	if LRU.list.Size() > LRU.capacity {
		delete(LRU.data, LRU.list.DeleteHead().GetKey())
	}
	LRU.data[key] = item
}
