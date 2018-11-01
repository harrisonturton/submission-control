package cache

import (
	"sync"
	"time"
)

// Each item inside the cache
type item struct {
	data      []byte
	timestamp time.Time
}

// Cache is a thread-safe store of result data.
type Cache struct {
	results           map[string]item
	mux               sync.Mutex
	opIndex           int
	evictionFrequency int
	evictionAge       time.Duration
}

// New creates a new Cache instance
func New(evictionFrequency int, evictionAge time.Duration) *Cache {
	return &Cache{
		results:           make(map[string]item),
		opIndex:           0,
		evictionFrequency: evictionFrequency,
		evictionAge:       evictionAge,
	}
}

// Set will set the value for a given string.
// If the value already exists, it will be overwritten.
func (c *Cache) Set(key string, value []byte) {
	c.mux.Lock()
	c.results[key] = item{value, time.Now()}
	c.opIndex++
	if c.opIndex%c.evictionFrequency == 0 {
		c.evict()
	}
	c.mux.Unlock()
}

// Get will fetch the data for a given string.
// If it doesn't exist, the second return value will
// be false.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	val, ok := c.results[key]
	c.opIndex++
	if c.opIndex%c.evictionFrequency == 0 {
		c.evict()
	}
	c.mux.Unlock()
	return val.data, ok
}

// Delete will remove an item from the cache.
// If it doesn't exist, then it is a no-op.
func (c *Cache) Delete(key string) {
	c.mux.Lock()
	delete(c.results, key)
	c.mux.Unlock()
}

// Length returns the number of items stored in
// the cache.
func (c *Cache) Length() int {
	return len(c.results)
}

// evict will run through the map and remove
// everything older than evictionAge
func (c *Cache) evict() {
	for key, value := range c.results {
		if time.Since(value.timestamp) > c.evictionAge {
			delete(c.results, key)
		}
	}
}
