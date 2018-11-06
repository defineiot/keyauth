package cache

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bluele/gcache"
)

type memCache struct {
	cache gcache.Cache
	size  int
}

// Newmemcache use to
func Newmemcache(size int) Cache {
	gc := gcache.New(size).LRU().Build()

	mcache := memCache{cache: gc}
	return &mcache
}

// Set use to set value to cache
func (m *memCache) Set(k string, x interface{}, ttl time.Duration) bool {
	b, err := json.Marshal(x)
	if err != nil {
		log.Printf("[Cache] E! marshal object error, %s\n", err)
		return false
	}

	if ttl == 0 {
		if err := m.cache.Set(k, b); err != nil {
			log.Printf("[Cache] E! set no ttl cache key error, %s\n", err)
			return false
		}
	} else {
		if err := m.cache.SetWithExpire(k, b, ttl); err != nil {
			log.Printf("[Cache] E! set ttl cache key error, %s\n", err)
			return false
		}
	}

	return true
}

// Get use to get value from cache
func (m *memCache) Get(k string, x interface{}) bool {
	v, err := m.cache.Get(k)
	if err != nil {
		log.Printf("[Cache] E! get value from cache error, %s\n", err)
		return false
	}

	log.Println(string(v.([]byte)))

	if err := json.Unmarshal(v.([]byte), x); err != nil {
		log.Printf("[Cache] E! ummarshal object error, %s\n", err)
		return false
	}

	return true
}

// Delete use to get value from cache
func (m *memCache) Delete(k string) bool {
	m.cache.Remove(k)
	return true
}
