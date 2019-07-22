package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCache(t *testing.T) {
	var value string
	cache := NewInMemoryCache()

	err := cache.Get("novalue", &value)
	assert.Equal(t, ErrNotFound, err)
	assert.Equal(t, "", value)

	err = cache.Set("key", "test-value", 0)
	assert.NoError(t, err)

	err = cache.Get("key", &value)
	assert.NoError(t, err)
	assert.Equal(t, "test-value", value)

	_ = cache.Del("key")
	err = cache.Get("key", &value)
	assert.Equal(t, ErrNotFound, err)
}

func TestMSET(t *testing.T) {
	cache := NewInMemoryCache()
	pairs := map[string]interface{}{
		"key1": "test value 1",
		"key2": "test value 2",
		"key3": "test value 3",
	}

	err := cache.MSet(pairs, 0)
	assert.NoError(t, err)

	for key, expected := range pairs {
		var actual string
		err = cache.Get(key, &actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestExpiration(t *testing.T) {
	var value string
	cache := NewInMemoryCache()
	pairs := map[string]interface{}{
		"key1": "test value 1",
		"key2": "test value 2",
		"key3": "test value 3",
	}

	err := cache.Set("key", "test-value", time.Millisecond*100)
	assert.NoError(t, err)

	err = cache.MSet(pairs, time.Millisecond*100)
	assert.NoError(t, err)

	err = cache.Get("key", &value)
	assert.NoError(t, err)
	assert.Equal(t, "test-value", value)

	for key, expected := range pairs {
		var actual string
		err = cache.Get(key, &actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	time.Sleep(time.Millisecond * 100)

	err = cache.Get("key", &value)
	assert.Equal(t, ErrNotFound, err)

	for key := range pairs {
		var actual string
		err = cache.Get(key, &actual)
		assert.Equal(t, ErrNotFound, err)
	}
}
