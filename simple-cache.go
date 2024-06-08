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

func NewSimpleCache(expiration time.Duration) *SimpleCache {
	c := &SimpleCache{
		expiration: expiration,
		items:      make(map[string]item),
	}
	go c.Cleaner(expiration)
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

func (c *SimpleCache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.rw.Lock()
	defer c.rw.Unlock()
	for k, v := range c.items {
		if now > v.expiration {
			delete(c.items, k)
		}
	}
}

func (c *SimpleCache) Cleaner(expiration time.Duration) {
	ticker := time.NewTicker(expiration)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		}
	}
}
