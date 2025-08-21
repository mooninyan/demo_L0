package mapcache

import (
	"demoL0/internal/utils"
	"sync"
)

type MapCache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

func NewCache[K comparable, V any]() *MapCache[K, V] {
	return &MapCache[K, V]{
		items: make(map[K]V),
	}
}

func (c *MapCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *MapCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.items[key]
	return val, ok
}

func (c *MapCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *MapCache[K, V]) GetAll() ([]V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.items) > 0 {
		values := utils.Values(c.items)
		return values, true
	}
	return nil, false
}
