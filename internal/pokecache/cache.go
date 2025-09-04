package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu sync.Mutex
	interval time.Duration
}

func (c *Cache) Add(key string, val []byte){
	c.mu.Lock()
	defer c.mu.Unlock()
	
	newEntry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	c.entries[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool){
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if ok{
		return entry.val, true
	}else{
		return nil, false
	}
}

func (c *Cache) reapLoop(){
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		<-ticker.C

		c.mu.Lock()
		now := time.Now()
		for k, entry := range c.entries{
			if now.Sub(entry.createdAt) > c.interval{
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache{
	c := &Cache{
		entries: make(map[string]cacheEntry),
		interval: interval,
	}

	go c.reapLoop()

	return c
}