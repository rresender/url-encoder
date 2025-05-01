package cache

import (
	"sync"

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
