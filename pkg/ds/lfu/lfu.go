package lfu

import (
	"homework/pkg/ds/linkedlist"
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	lock sync.RWMutex

	nodeMap  map[K]*item[K, V]
	listMap  map[int]*linkedlist.List[K, V]
	capacity int
	min      int

	ttl time.Duration
}

func NewLFU[K comparable, V any](capacity int, ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		nodeMap:  make(map[K]*item[K, V]),
		listMap:  make(map[int]*linkedlist.List[K, V]),
		capacity: capacity,
		min:      0,
		ttl:      ttl,
	}
}

func (l *Cache[K, V]) Get(key K) (V, bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.get(key)
}

func (l *Cache[K, V]) Remove(key K) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	node, ok := l.nodeMap[key]
	if !ok {
		return false
	}
	delete(l.nodeMap, key)

	list, ok := l.listMap[node.freq]
	if ok {
		list.DeleteNode(node.node)
	}

	return true
}

func (l *Cache[K, V]) get(key K) (V, bool) {
	node, ok := l.nodeMap[key]
	if !ok {
		return defaultValue[V](), false
	}

	list, ok := l.listMap[node.freq]
	if ok {
		list.DeleteNode(node.node)
	}
	if node.Expired(time.Now()) {
		delete(l.nodeMap, node.node.GetKey())
		return defaultValue[V](), false
	}

	node.freq++
	nextList, nextOk := l.listMap[node.freq]
	if !nextOk {
		nextList = linkedlist.New[K, V]()
	}
	nextList.PushNode(node.node)
	l.listMap[node.freq] = nextList
	if list.Size() == 0 && l.min == node.freq-1 {
		l.min++
	}
	return node.node.GetValue(), true
}

func (l *Cache[K, V]) Put(key K, value V) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.capacity == 0 {
		return
	}

	node, ok := l.nodeMap[key]
	if ok {
		node.node.SetValue(value)
		l.get(key)
		return
	}

	if len(l.nodeMap) == l.capacity {
		minList := l.listMap[l.min]
		node := minList.DeleteHead()
		delete(l.nodeMap, node.GetKey())
	}

	node = &item[K, V]{
		node:      linkedlist.NewNode[K, V](key, value),
		freq:      1,
		expiredAt: time.Now().Add(l.ttl),
	}

	l.min = 1
	list, ok := l.listMap[node.freq]
	if !ok {
		list = linkedlist.New[K, V]()
	}

	list.PushNode(node.node)
	l.listMap[node.freq] = list
	l.nodeMap[key] = node
}

func defaultValue[V any]() V {
	var v V
	return v
}
