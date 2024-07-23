package cache

import (
	"github.com/stretchr/testify/require"
	"homework/internal/model"
	"testing"
	"time"
)

var (
	key    = "key1"
	order1 = model.Order{ID: "1"}
	order2 = model.Order{ID: "2"}
	value  = []model.Order{
		order1,
		order2,
	}
)

func TestCacheOrders_Empty(t *testing.T) {
	cache := NewOrdersCache(10, time.Hour)
	_, ok := cache.Get(key)
	require.False(t, ok)
}

func TestCacheOrders_PutGet(t *testing.T) {
	cache := NewOrdersCache(10, time.Hour)
	cache.Put(key, value)

	cached, ok := cache.Get(key)

	require.True(t, ok)
	require.Equal(t, value, cached)
}

func TestCacheOrders_PutRemoveByID(t *testing.T) {
	cache := NewOrdersCache(10, time.Hour)
	cache.Put(key, value)

	cache.RemoveById(order1.ID)

	_, ok := cache.Get(key)
	require.False(t, ok)
}

func TestCacheOrders_PutRemoveByIDs(t *testing.T) {
	cache := NewOrdersCache(10, time.Hour)
	cache.Put(key, value)

	cache.RemoveByIds([]string{order1.ID, order2.ID})

	_, ok := cache.Get(key)
	require.False(t, ok)
}
