package pokecache

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add("Test", []byte("This is test data"))
}

func TestGet(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add("Hiya", []byte{})
	if _, ok := cache.Get("Hiya"); !ok {
		t.Fatal("Reap loop removed the entry before time was up.")
	}
}

func TestReapLoop(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add("test", []byte{})
	if _, ok := cache.Get("test"); !ok {
		t.Fatal("Reap loop removed the entry before time was up.")
	}
	time.Sleep(5 * time.Second)
	if _, ok := cache.Get("test"); ok {
		t.Fatalf("Reap loop failed to remove test entry")
	}
}
