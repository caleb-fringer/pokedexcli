package pokecache

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"
)

var testUrl *url.URL

func init() {
	var err error
	testUrl, err = url.Parse("http://test.go/path/to/endpoint")
	if err != nil {
		fmt.Println("Error initializing test url.")
		os.Exit(1)
	}
}

func TestAdd(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add(*testUrl, []byte("This is test data"))
}

func TestGet(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add(*testUrl, []byte{})
	time.Sleep(time.Second)
	if _, ok := cache.Get(*testUrl); !ok {
		t.Fatal("Reap loop removed the entry before time was up.")
	}
}

func TestReapLoop(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add(*testUrl, []byte{})
	if _, ok := cache.Get(*testUrl); !ok {
		t.Fatal("Reap loop removed the entry before time was up.")
	}
	time.Sleep(5 * time.Second)
	if _, ok := cache.Get(*testUrl); ok {
		t.Fatalf("Reap loop failed to remove *test entry")
	}
}
