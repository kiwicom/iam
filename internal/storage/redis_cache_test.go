package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheKey(t *testing.T) {
	redis := RedisCache{version: 1}

	key := redis.cacheKey("user")

	assert.Equal(t, "user-v1", key)
}
