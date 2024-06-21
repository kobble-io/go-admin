package utils

import (
	"sync"
	"time"
)

type CacheKey[T any] struct {
	Data      T
	CreatedAt int64
	ExpiresAt *int64
}

type CacheConfig struct {
	DefaultTtl *time.Duration
}

type Cache[T any] struct {
	data   map[string]CacheKey[T]
	config CacheConfig
	mu     sync.Mutex
}

func NewCache[T any](config CacheConfig) *Cache[T] {
	return &Cache[T]{
		data:   make(map[string]CacheKey[T]),
		config: config,
	}
}

func (c *Cache[T]) Get(key string) *T {
	c.mu.Lock()
	defer c.mu.Unlock()

	ckey, exists := c.data[key]
	if !exists {
		return nil
	}

	now := time.Now().UnixMilli()
	if ckey.ExpiresAt != nil && now > *ckey.ExpiresAt {
		delete(c.data, key)
		return nil
	}

	return &ckey.Data
}

func (c *Cache[T]) Set(key string, data T, ttl *time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	var expiresAt *int64
	if ttl != nil {
		expires := now.Add(*ttl * time.Second).UnixMilli()
		expiresAt = &expires
	} else if c.config.DefaultTtl != nil {
		expires := now.Add(*c.config.DefaultTtl).UnixMilli()
		expiresAt = &expires
	}

	c.data[key] = CacheKey[T]{
		Data:      data,
		CreatedAt: now.UnixMilli(),
		ExpiresAt: expiresAt,
	}
}
