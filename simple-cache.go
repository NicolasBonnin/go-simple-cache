package go_simple_cache

import (
	"sync"
	"time"
)

type SimpleCache struct {
	items      map[string]item
	rw         sync.RWMutex
	expiration time.Duration
}

type item struct {
	object     interface{}
	expiration int64
}

// NewSimpleCache initializes a new SimpleCache with a default expiration time of 1 hour.
// Optionally, you can set a custom expiration time by passing a duration parameter.
func NewSimpleCache() *SimpleCache {
	items := make(map[string]item)
	c := &SimpleCache{
		items:      items,
		expiration: 1 * time.Hour, // Set default expiration time.
	}
	return c
}

func (c *SimpleCache) Set(key string, object interface{}) {
	c.rw.Lock()
	var e = time.Now().Add(c.expiration).UnixNano()
	c.items[key] = item{
		object:     object,
		expiration: e,
	}
	c.rw.Unlock()
}

func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.rw.RLock()
	item, found := c.items[key]
	if !found {
		c.rw.RUnlock()
		return nil, false
	}
	c.rw.RUnlock()
	if item.expired() {
		c.Delete(key)
		return nil, false
	}
	return item.object, true
}

func (item item) expired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.expiration
}

func (c *SimpleCache) Delete(key string) {
	c.rw.Lock()
	delete(c.items, key)
	c.rw.Unlock()
}

func (c *SimpleCache) Flush() {
	c.rw.Lock()
	c.items = map[string]item{}
	c.rw.Unlock()
}
