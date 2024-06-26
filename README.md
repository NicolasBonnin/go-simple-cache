# Go Simple Cache
Go Simple Cache is a lightweight, thread-safe, and easy-to-use caching library for Go. It provides a simple key-value store with expiration support, ideal for caching frequently accessed data to improve performance.

## Features
 * **Simple Interface**: Straightforward API for setting, getting, deleting, and flushing items in the cache.
 * **Expiration Support**: Set expiration times for cached items to automatically remove stale data.
 * **Thread-Safe**: Utilizes mutex locks for safe concurrent access to the cache.
 * **Lightweight**: Minimal overhead and dependencies, making it suitable for a wide range of applications.
 * **Easy Integration**: Drop-in library that can be quickly integrated into existing projects.


## Installation
You can install Go Simple Cache using go get:

`go get github.com/nicolasbonnin/go-simple-cache`

## Usage

```go
    package main

    import (
        "fmt"
        "time"
        "github.com/nicolasbonnin/go-simple-cache"
    )

    func main() {
        // Create a new instance of the cache
        cache := go_simple_cache.NewSimpleCache(1 * time.Hour)
	    
        // Set a key-value pair in the cache
        cache.Set("key", "value")
	    
        // Retrieve a value from the cache
        val, found := cache.Get("key")
        if found {
        	fmt.Println("Value found:", val)
        } else {
        	fmt.Println("Value not found")
        }
	
        // Delete an item from the cache
        cache.Delete("key")
	
        // Flush the entire cache
        cache.Flush()
    }
```

## Use Case: Caching Database Results

Example of how to integrate the cache with database queries:
```go
    // FindByNameWithCache fetches data by name from the cache if available,
    // otherwise, it fetches it from the database, caches it, and returns it.
    func (s *Service) FindByNameWithCache(name string) (*Data, error) {
        data, found := s.getDataFromCache(name)
        if !found {
            // Fetch data from the database
            data, err := s.FindByName(name)
            if err != nil {
                return nil, err
            }
            // Cache the fetched data
            s.simpleCache.Set(name, data)
            return data, nil
        }

        return data, nil
    }

    // CacheFlushByName removes the cached entry for the given name.
    func (s *Service) CacheFlushByName(name string) {
        s.simpleCache.Delete(name)
    }
	
    // getDataFromCache retrieves data from the cache by name.
    func (s *Service) getDataFromCache(name string) (*Data, bool) {
        cachedData, found := s.simpleCache.Get(name)
        if found {
            data, ok := cachedData.(*Data)
			
            return data, ok
        }
        return nil, false
    }
```