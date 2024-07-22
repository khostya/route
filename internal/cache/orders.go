package cache

import (
	"homework/internal/model"
	"homework/pkg/ds/lfu"
	"sync"
	"time"
)

type (
	KeyOrder = string

	OrdersCache struct {
		cache Cache[KeyOrder, []model.Order]

		getKeyByID map[string][]KeyOrder

		lock sync.RWMutex
	}
)

func NewOrdersCache(cap int, ttl time.Duration) *OrdersCache {
	return &OrdersCache{
		cache:      lfu.NewLFU[KeyOrder, []model.Order](cap, ttl),
		getKeyByID: make(map[string][]KeyOrder),
	}
}

func (o *OrdersCache) Put(k string, v []model.Order) {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.cache.Put(k, v)

	for _, order := range v {
		_, ok := o.getKeyByID[order.ID]
		if !ok {
			o.getKeyByID[order.ID] = make([]KeyOrder, 0, 1)
		}
		o.getKeyByID[order.ID] = append(o.getKeyByID[order.ID], k)
	}
}

func (o *OrdersCache) RemoveById(id string) {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.removeById(id)
}

func (o *OrdersCache) removeById(id string) {
	keys, ok := o.getKeyByID[id]
	if !ok {
		return
	}

	for _, k := range keys {
		o.cache.Remove(k)
	}
	delete(o.getKeyByID, id)
}

func (o *OrdersCache) RemoveByIds(ids []string) {
	o.lock.Lock()
	defer o.lock.Unlock()

	for _, id := range ids {
		o.removeById(id)
	}
}

func (o *OrdersCache) Get(k string) ([]model.Order, bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	return o.cache.Get(k)
}
