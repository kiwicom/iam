package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

type mockCache struct {
	mock.Mock
	data map[string]cacheItem
}

// `value` is overwritten before use
// nolint: staticcheck
func (c *mockCache) Get(key string, value interface{}) error {
	v, ok := c.data[key]
	if !ok || time.Since(v.expiration) > 0 {
		delete(c.data, key)
		return ErrNotFound
	}
	// ineffective assignment of `value`
	// nolint: ineffassign
	value = &v
	return nil
}
func (c *mockCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
	return nil
}
func (c *mockCache) Del(key string) error {
	delete(c.data, key)
	return nil
}

func newMocks() (*mockCache, *LockManager) {
	cache := &mockCache{
		data: make(map[string]cacheItem),
	}
	lock := &LockManager{
		cache:      cache,
		retryDelay: 10 * time.Millisecond,
		expiration: 100 * time.Millisecond,
	}

	return cache, lock
}

func TestCreate(t *testing.T) {
	cache, lock := newMocks()

	// Create lock successfully
	err := lock.Create("test")
	assert.Nil(t, err)
	_, ok := cache.data["lock:test"]
	assert.Equal(t, true, ok)

	// Creating a second lock should throw an error, because we already have an
	// existing lock.
	err = lock.Create("test")
	assert.Equal(t, ErrLockExists, err)

	// After failing to create a lock, there should be no lock left. Since
	// lock.Create should wait for the existing lock to be deleted or expired.
	_, ok = cache.data["lock:test"]
	assert.Equal(t, false, ok)
}

func TestDelete(t *testing.T) {
	cache, lock := newMocks()

	// Create lock successfully
	err := lock.Create("test")
	assert.Nil(t, err)
	_, ok := cache.data["lock:test"]
	assert.Equal(t, true, ok)

	// Delete lock successfully
	lock.Delete("test")
	_, ok = cache.data["lock:test"]
	assert.Equal(t, false, ok)
}
