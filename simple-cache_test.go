package go_simple_cache

import (
	"sync"
	"testing"
	"time"
)

func TestExpiredGetLeavesCleanupToCleaner(t *testing.T) {
	c := &SimpleCache{items: map[string]item{
		"key": {object: "expired", expiration: time.Now().Add(-time.Hour).UnixNano()},
	}}
	if value, found := c.Get("key"); found || value != nil {
		t.Fatalf("Get expired entry = (%v, %v), want (nil, false)", value, found)
	}
	if _, exists := c.items["key"]; !exists {
		t.Fatal("Get removed an entry instead of leaving cleanup to the cleaner")
	}
	c.DeleteExpired()
	if _, exists := c.items["key"]; exists {
		t.Fatal("cleaner did not remove the expired entry")
	}
}

func TestCacheOperations(t *testing.T) {
	c := &SimpleCache{items: make(map[string]item), expiration: time.Hour}
	if value, found := c.Get("missing"); found || value != nil {
		t.Fatalf("Get missing entry = (%v, %v), want (nil, false)", value, found)
	}
	c.Set("key", "old")
	c.Set("key", "fresh")
	c.DeleteExpired()
	if value, found := c.Get("key"); !found || value != "fresh" {
		t.Fatalf("Get after replacement and cleanup = (%v, %v), want (fresh, true)", value, found)
	}
	c.Delete("key")
	if _, found := c.Get("key"); found {
		t.Fatal("deleted entry was found")
	}
	c.Set("first", 1)
	c.Set("second", 2)
	c.Flush()
	if len(c.items) != 0 {
		t.Fatal("Flush left entries in the cache")
	}
}

func TestConcurrentCacheOperations(t *testing.T) {
	c := &SimpleCache{items: make(map[string]item), expiration: time.Hour}
	var workers sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for i := 0; i < 100; i++ {
				c.Set("key", i)
				c.Get("key")
				c.DeleteExpired()
				c.Delete("key")
				c.Flush()
			}
		}()
	}
	workers.Wait()
	c.Set("key", "fresh")
	if value, found := c.Get("key"); !found || value != "fresh" {
		t.Fatalf("Get final entry = (%v, %v), want (fresh, true)", value, found)
	}
}
