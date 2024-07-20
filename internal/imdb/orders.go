package imdb

import (
	"homework/internal/model"
	"homework/pkg/ds/lfu"
	"sync"
	"time"
)

type (
	KeyOrder = string

	OrdersIMDB struct {
		imdb IMDB[KeyOrder, []model.Order]

		getKeyByID map[string][]KeyOrder

		lock sync.RWMutex
	}
)

func NewOrdersIMDB(cap int, ttl time.Duration) *OrdersIMDB {
	return &OrdersIMDB{
		imdb:       lfu.NewLFU[KeyOrder, []model.Order](cap, ttl),
		getKeyByID: make(map[string][]KeyOrder),
	}
}

func (o *OrdersIMDB) Put(k string, v []model.Order) {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.imdb.Put(k, v)

	for _, order := range v {
		_, ok := o.getKeyByID[order.ID]
		if !ok {
			o.getKeyByID[order.ID] = make([]KeyOrder, 0, 1)
		}
		o.getKeyByID[order.ID] = append(o.getKeyByID[order.ID], k)
	}
}

func (o *OrdersIMDB) RemoveById(id string) bool {
	o.lock.Lock()
	defer o.lock.Unlock()
	keys, ok := o.getKeyByID[id]
	if !ok {
		return false
	}

	for _, k := range keys {
		o.imdb.Remove(k)
	}
	return true
}

func (o *OrdersIMDB) Get(k string) ([]model.Order, bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	return o.imdb.Get(k)
}
