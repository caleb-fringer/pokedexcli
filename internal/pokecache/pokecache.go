package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	mutex    *sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) (cache Cache) {
	cache = Cache{
		entries:  make(map[string]cacheEntry),
		mutex:    &sync.Mutex{},
		interval: interval,
	}
	go cache.reapLoop()
	return cache
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

func (cache *Cache) reapLoop() {
	ticker := time.NewTicker(cache.interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		cache.mutex.Lock()
		for key, val := range cache.entries {
			if time.Since(val.createdAt) >= cache.interval {
				delete(cache.entries, key)
			}
		}
		cache.mutex.Unlock()
	}
}
