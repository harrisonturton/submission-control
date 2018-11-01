package cache

import (
	"testing"
	"time"
)

var evictionFrequency = 10      // Do an eviction check after this number of operations
var evictionAge = 5 * time.Hour // Evict items that are older than this

func TestAdd(t *testing.T) {
	cache := New(evictionFrequency, evictionAge)
	key := "key"
	val := []byte("val")
	cache.Set(key, val)
	if cache.Length() != 1 {
		t.Errorf("Failed to add. Expected length %d, but got %d", 0, cache.Length())
	}
}

func TestRemove(t *testing.T) {
	cache := New(evictionFrequency, evictionAge)

	key := "key"
	val := []byte("val")
	cache.Set(key, val)
	cache.Delete(key)
	if cache.Length() != 0 {
		t.Errorf("Failed to delete. Expected length %d, got %d", 0, cache.Length())
	}
}

func TestLength(t *testing.T) {
	cache := New(evictionFrequency, evictionAge)
	checkLength := func(length int) {
		if cache.Length() != length {
			t.Errorf("Expected length %d, but got %d", length, cache.Length())
		}
	}
	checkLength(0)
	cache.Set("one", []byte(""))
	checkLength(1)
	cache.Set("two", []byte(""))
	checkLength(2)
	for i := 0; i < 100; i++ {
		cache.Set(string(i), []byte("data"))
	}
	checkLength(102)
	cache.Delete("two")
	checkLength(101)
	cache.Delete("one")
	checkLength(100)
	for i := 0; i < 100; i++ {
		cache.Delete(string(i))
	}
	checkLength(0)
}
