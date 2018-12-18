package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defineiot/keyauth/internal/cache"
)

func TestRedisSet(t *testing.T) {
	t.Run("RedisSetOK", testRedisSetOK)
}

func testRedisSetOK(t *testing.T) {
	var tk token

	rc, err := cache.NewRedisCache("172.168.0.90:6379", 1, "")
	assert.NoError(t, err)

	rc.Set("token", intoken, 0)

	rc.Get("token", &tk)
	assert.Equal(t, intoken, tk)
}

func TestRedisDel(t *testing.T) {
	t.Run("RedisSetOK", testRedisDelOK)
}

func testRedisDelOK(t *testing.T) {
	var tk token

	rc, err := cache.NewRedisCache("172.168.0.90:6379", 1, "")
	assert.NoError(t, err)

	rc.Set("token", intoken, 0)

	ok := rc.Delete("token")
	assert.Equal(t, true, ok)

	ok = rc.Get("token", &tk)
	assert.Equal(t, false, ok)

}
