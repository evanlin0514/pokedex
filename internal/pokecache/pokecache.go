package pokecache

import (
 	"sync"
	"time"
)

type cache struct {
	data map[string]cacheEntry
	mu *sync.RMutex
}

type cacheEntry struct{
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) *cache{
	c := &cache{
		data: make(map[string]cacheEntry)
	}

	go func() {
		for {
			time.Sleep(internal)
			
		}
	}
}