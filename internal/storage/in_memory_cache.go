package storage

import (
	"strings"
	"time"
)

// item represent a value in cache
type item struct {
	value      []byte
	expiration time.Time
}

// InMemoryCache is an in memory cache used as a backup when Redis is unavailable
type InMemoryCache map[string]item

// NewInMemoryCache initializes and returns an InMemoryCache
func NewInMemoryCache() InMemoryCache {
	return make(map[string]item)
}

// Get retrieves an item from cache.
// `key` is case insensitive.
// `value` is a pointer to the variable that will receive the data.
// `error` is ErrNotFound when no value is found
func (c InMemoryCache) Get(key string, value interface{}) error {
	lowerKey := strings.ToLower(key)
	data, ok := c[lowerKey]
	if ok {
		if time.Now().Before(data.expiration) || data.expiration.IsZero() {
			return json.Unmarshal(data.value, &value)
		}
		// Item is expired
		_ = c.Del(key)
	}

	return ErrNotFound
}

// Set writes data to cache. `key` is case insensitive.
func (c InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	strVal, err := json.Marshal(value)
	if err != nil {
		return err
	}

	expiration := time.Time{}
	if ttl != 0 {
		expiration = time.Now().Add(ttl)
	}

	lowerKey := strings.ToLower(key)
	c[lowerKey] = item{
		strVal,
		expiration,
	}
	return nil
}

// Del deletes an item from cache
func (c InMemoryCache) Del(key string) error {
	lowerKey := strings.ToLower(key)
	delete(c, lowerKey)
	return nil
}

// MSet writes items to cache in bulk
func (c InMemoryCache) MSet(pairs map[string]interface{}, ttl time.Duration) error {
	bytePairs := make(map[string][]byte)

	// Go through all values and convert them to byte arrays first, then write to
	// cache. This is kept deliberately in two separate steps to ensure that if
	// there is an error nothing is written in cache.
	for key, value := range pairs {
		strValue, err := json.Marshal(value)
		if err != nil {
			return err
		}
		bytePairs[key] = strValue
	}

	expiration := time.Time{}
	if ttl != 0 {
		expiration = time.Now().Add(ttl)
	}
	for key, value := range bytePairs {
		lowerKey := strings.ToLower(key)
		c[lowerKey] = item{
			value,
			expiration,
		}
	}
	return nil
}
