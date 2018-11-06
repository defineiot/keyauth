package cache

import "time"

// Cache provides the interface for cache implementations.
type Cache interface {
	//Set stores a value with the given ttl
	Set(key string, value interface{}, ttl time.Duration) bool

	//Get retrieves a value previously stored in the cache.
	//value has to be a pointer to a data structure that matches the type previously given to Set
	//The return value indicates if a value was found
	Get(key string, value interface{}) bool

	// Delete cache
	Delete(key string) bool
}
