package inmem

import (
	"context"
	"sync"
	"time"
)

type cacheEntry struct {
	value     string
	expiresAt time.Time
}

// Cache is a thread-safe in-memory implementation of health.CachePort.
// It is used in tests and local development where a real Valkey instance is not available.
type Cache struct {
	mu    sync.RWMutex
	store map[string]cacheEntry
}

// NewCache returns an empty in-memory Cache ready for use.
func NewCache() *Cache {
	return &Cache{store: make(map[string]cacheEntry)}
}

func (c *Cache) Get(_ context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.store[key]
	if !ok || time.Now().After(e.expiresAt) {
		return "", nil
	}
	return e.value, nil
}

func (c *Cache) Set(_ context.Context, key, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = cacheEntry{value: value, expiresAt: time.Now().Add(ttl)}
	return nil
}

func (c *Cache) Del(_ context.Context, keys ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, k := range keys {
		delete(c.store, k)
	}
	return nil
}
