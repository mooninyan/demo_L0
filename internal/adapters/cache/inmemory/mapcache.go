package inmemory

import (
	"sync"
	"time"
)

// cacheItem хранит value и его expiration time.
type cacheItem[V any] struct {
	value      V
	expiration int64
}

// MapCache потокобезопасный in-memory cache с TTL инвалидацией.
type MapCache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]cacheItem[V]
	ttl   time.Duration
}

// NewCache создает MapCache с TTL и интервалом очистки.
// Стартует горутину очищающую кэш от истекших по TTL элементов.
func NewCache[K comparable, V any](defaultTTL time.Duration, cleanupInterval time.Duration) *MapCache[K, V] {
	c := &MapCache[K, V]{
		items: make(map[K]cacheItem[V]),
		ttl:   defaultTTL,
	}

	if cleanupInterval > 0 {
		go c.startCleanup(cleanupInterval)
	}

	return c
}

// Set добавляет элемент в кэш со значением TTL по умолчанию.
func (c *MapCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiration := time.Now().Add(c.ttl).UnixNano()
	c.items[key] = cacheItem[V]{
		value:      value,
		expiration: expiration,
	}
}

// Get Извлекает элемент из кэша. Возвращает значение и значение true, если элемент существует и срок его действия не истек.
// В противном случае возвращает нулевое значение для типа и значение false.
func (c *MapCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok {
		var zeroV V
		return zeroV, false
	}

	if time.Now().UnixNano() > item.expiration {
		var zeroV V
		return zeroV, false // Expired
	}

	return item.value, true
}

// Delete удаляет элемент из кэша.
func (c *MapCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// GetAll извлекает все непросроченные элементы из кэша.
func (c *MapCache[K, V]) GetAll() ([]V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.items) == 0 {
		return nil, false
	}

	values := make([]V, 0, len(c.items))
	now := time.Now().UnixNano()
	for _, item := range c.items {
		if now < item.expiration {
			values = append(values, item.value)
		}
	}

	if len(values) > 0 {
		return values, true
	}

	return nil, false
}

// startCleanup запускает цикл, который запускает очистку кэша с заданным интервалом.
func (c *MapCache[K, V]) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.deleteExpired()
	}
}

// deleteExpired удаляет все просроченные элементы из кэша.
func (c *MapCache[K, V]) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if now > v.expiration {
			delete(c.items, k)
		}
	}
}
