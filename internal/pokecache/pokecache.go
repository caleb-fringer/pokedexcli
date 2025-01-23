package pokecache

import (
	"net/url"
	"sync"
	"time"
)

type Cache struct {
	sync.Mutex
	entries  map[url.URL]cacheEntry
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) (cache *Cache) {
	cache = &Cache{
		entries:  make(map[url.URL]cacheEntry),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (cache *Cache) Add(key url.URL, val []byte) {
	entry := cacheEntry{time.Now(), val}
	cache.Lock()
	defer cache.Unlock()
	cache.entries[key] = entry
}

func (cache *Cache) Get(key url.URL) (entryData []byte, ok bool) {
	cache.Lock()
	defer cache.Unlock()
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
		cache.Lock()
		for key, val := range cache.entries {
			if time.Since(val.createdAt) >= cache.interval {
				delete(cache.entries, key)
			}
		}
		cache.Unlock()
	}
}
