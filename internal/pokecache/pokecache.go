package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entry map[string]cacheEntry
	mu *sync.Mutex
}

type cacheEntry struct{
	createdAt time.Time
	Val []byte
}

func NewCache(interval time.Duration) *Cache{
	c := &Cache{
		Entry: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}


	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
				case <- ticker.C:
					c.reapLoop(interval)	
			}
		}
	}()
	return c
}

func (c *Cache) Add(key string, val []byte){
	c.mu.Lock()
	defer c.mu.Unlock() 

	c.Entry[key] = cacheEntry{
		createdAt: time.Now(),
		Val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if data, ok := c.Entry[key]; ok {
		return data.Val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//loop throught caches to find expired entry
	now := time.Now()
	for key, entry := range c.Entry{
		if entry.createdAt.Add(interval).Before(now){
			delete(c.Entry, key)
		} 
	}
}