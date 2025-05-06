package cache

import (
	"sync"

	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/rresender/url-enconder/internal/model"
)

type EncodeURLCache interface {
	Set(key string, value *model.EncodeURL)
	Get(key string) (*model.EncodeURL, bool)
}

type inMemoryCache struct {
	cache map[string]*model.EncodeURL
	mutex sync.RWMutex
}

func NewInMemoryCache() EncodeURLCache {
	return &inMemoryCache{
		cache: make(map[string]*model.EncodeURL),
	}
}

func (c *inMemoryCache) Set(key string, value *model.EncodeURL) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = value
}

func (c *inMemoryCache) Get(key string) (*model.EncodeURL, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	val, ok := c.cache[key]
	return val, ok
}

func NewInMemoryTTLCache(cacheTTL time.Duration) EncodeURLCache {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, *model.EncodeURL](cacheTTL),
	)
	go cache.Start()
	return &ttlCache{cache: cache}
}

type ttlCache struct {
	cache *ttlcache.Cache[string, *model.EncodeURL]
}

func (c *ttlCache) Set(key string, value *model.EncodeURL) {
	c.cache.Set(key, value, ttlcache.DefaultTTL)
}

func (c *ttlCache) Get(key string) (*model.EncodeURL, bool) {
	if c.cache.Has(key) {
		return c.cache.Get(key).Value(), true
	}
	return nil, false
}
