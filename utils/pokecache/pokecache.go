package pokecache

import (
	"sync"
	"time"
)

type cache struct {
	entries map[string]cacheEntry
	mutex   sync.Mutex
}

// cacheEntry is a struct that holds the creation time of the entry and the value

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// NewCache returns a new cache struct with a modifi
func NewCache(ttl time.Duration) *cache {
	c := &cache{
		entries: map[string]cacheEntry{},
	}
	go c.reapLoop(ttl)
	return c
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *cache) UpdateTime(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return
	}
	entry.createdAt = time.Now()
	c.entries[key] = entry
}

// reapLoop is a function that runs in the background (time.Ticker) and removes entries that are older than given time.Duration
func (c *cache) reapLoop(ttl time.Duration) {
	ticker := time.NewTicker(ttl)
	for range ticker.C {
		// Loop over entries and remove the ones that are older than ttl
		c.mutex.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > ttl {
				delete(c.entries, key)
			}
		}

		c.mutex.Unlock()
	}
}
