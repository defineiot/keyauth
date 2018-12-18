package cache

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type redisCache struct {
	client *redis.Client
}

// NewRedisCache todo
func NewRedisCache(address string, db int, password string) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &redisCache{client: client}, nil
}

func (r *redisCache) Set(key string, value interface{}, ttl time.Duration) bool {
	b, err := json.Marshal(value)
	if err != nil {
		log.Printf("[Cache] E! marshal object error, %s\n", err)
		return false
	}

	if err := r.client.Set(key, b, ttl).Err(); err != nil {
		log.Printf("[Cache] E! set cache error, %s\n", err)
		return false
	}

	return true
}

func (r *redisCache) Get(key string, value interface{}) bool {
	v, err := r.client.Get(key).Bytes()
	if err != nil {
		log.Printf("[Cache] E! get value from cache error, %s\n", err)
		return false
	}

	if err := json.Unmarshal(v, value); err != nil {
		log.Printf("[Cache] E! ummarshal object error, %s\n", err)
		return false
	}

	return true
}

func (r *redisCache) Delete(key string) bool {
	if err := r.client.Del(key).Err(); err != nil {
		log.Printf("[Cache] E! delete cache error, %s\n", err)
		return false
	}

	return true
}
