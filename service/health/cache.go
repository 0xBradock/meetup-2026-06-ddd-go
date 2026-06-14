package health

import (
	"context"
	"time"
)

// CachePort is the caching abstraction for the health domain.
// Implementations live in db/valkey (production) and db/inmem (tests).
type CachePort interface {
	// Get returns the cached value for key, or ("", nil) on a cache miss.
	Get(ctx context.Context, key string) (string, error)
	// Set stores value under key with the given TTL.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	// Del removes the given keys.
	Del(ctx context.Context, keys ...string) error
}
