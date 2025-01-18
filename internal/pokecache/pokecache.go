package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	mutex    sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	return Cache{entries: make(map[string]cacheEntry), interval: interval}
}

func (cache *Cache) Add(key string, val []byte) {
	entry := cacheEntry{time.Now(), val}
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.entries[key] = entry
}

func (cache *Cache) Get(key string) (entryData []byte, ok bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	entry, ok := cache.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}
