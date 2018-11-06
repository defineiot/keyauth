package cache_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/defineiot/keyauth/internal/cache"
)

var intoken = token{UUID: "aa-bb-cc", Token: "abc", Count: 10, Extra: map[string]string{"key": "value"}}

type token struct {
	UUID  string            `json:"uuid"`
	Token string            `json:"token"`
	Count int               `json:"count"`
	Extra map[string]string `json:"extra"`
}

func TestSet(t *testing.T) {
	t.Run("SetNoTTL", testSetNoTTL)
	t.Run("SetWithTTL", testSetWithTTL)
}

func testSetNoTTL(t *testing.T) {
	var tk token

	c := cache.Newmemcache(20)
	c.Set("token", intoken, 0)

	c.Get("token", &tk)
	assert.Equal(t, intoken, tk)
}

func testSetWithTTL(t *testing.T) {
	var tk token

	c := cache.Newmemcache(20)
	c.Set("token", intoken, time.Second*1)

	time.Sleep(time.Second * 1)
	ok := c.Get("token", &tk)
	assert.Equal(t, false, ok)
}

func TestDelete(t *testing.T) {
	t.Run("DeleteOK", testDeleteOK)
}

func testDeleteOK(t *testing.T) {
	var tk token

	c := cache.Newmemcache(20)
	c.Set("token", intoken, 0)

	ok := c.Delete("token")
	assert.Equal(t, true, ok)

	ok = c.Get("token", &tk)
	assert.Equal(t, false, ok)
}
