package storage

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	redisTrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// RedisCache contains redis client
type RedisCache struct {
	client  *redisTrace.Client
	backup  InMemoryCache
	version int
}

// ErrNotFound is returned when an item is not present in cache
var ErrNotFound = errors.New("item not found")

// shouldUseRedisBackup returns whether or not the inMemory backup for Redis
// should be used, depending on the error returned by Redis.
func shouldUseRedisBackup(err error) bool {
	return err != nil && err != redis.Nil
}

// NewRedisCache initializes and returns a RedisCache
func NewRedisCache(host, port string, version int) *RedisCache {
	opts := &redis.Options{Addr: net.JoinHostPort(host, port)}

	return &RedisCache{
		client: redisTrace.NewClient(opts, redisTrace.WithServiceName("kiwi-iam.redis")),
		backup: NewInMemoryCache(),
	}
}

// Get retrieves an item from cache.
// `key` is case insensitive.
// `value` is a pointer to the variable that will receive the data.
// `error` is ErrNotFound when no value is found.
func (c *RedisCache) Get(key string, value interface{}) error {
	lowerKey := c.cacheKey(key)
	data, err := c.client.Get(lowerKey).Bytes()
	if err == nil {
		err = json.Unmarshal(data, &value)
		return err
	}
	if err == redis.Nil {
		return ErrNotFound
	}

	if shouldUseRedisBackup(err) {
		log.Println("Redis down using inMemory GET")
		raven.CaptureMessage("Redis down using inMemory GET", nil)
		err = c.backup.Get(key, value)
	}
	return err
}

// Set writes data to cache with the specified lifespan
// `key` is case insensitive.
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	strVal, err := json.Marshal(value)
	if err != nil {
		return err
	}

	lowerKey := c.cacheKey(key)
	_, err = c.client.Set(lowerKey, strVal, ttl).Result()
	if shouldUseRedisBackup(err) {
		log.Println("Redis down using inMemory SET")
		raven.CaptureMessage("Redis down using inMemory SET", nil)
		err = c.backup.Set(key, value, ttl)
	}
	return err
}

// Del deletes an item from cache
func (c *RedisCache) Del(key string) error {
	lowerKey := c.cacheKey(key)
	_, err := c.client.Del(lowerKey).Result()
	if shouldUseRedisBackup(err) {
		log.Println("Redis down using inMemory DEL")
		raven.CaptureMessage("Redis down using inMemory DEL", nil)
		_ = c.backup.Del(key)
	}
	return err
}

// MSet writes items to cache in bulk
func (c *RedisCache) MSet(pairs map[string]interface{}, ttl time.Duration) error {
	args := make([]interface{}, len(pairs)*2)
	i := 0

	for key, value := range pairs {
		args[i] = c.cacheKey(key)
		strValue, err := json.Marshal(value)
		if err != nil {
			return err
		}
		args[i+1] = strValue
		i += 2
	}

	_, err := c.client.MSet(args...).Result()
	if shouldUseRedisBackup(err) {
		log.Println("Redis down using inMemory MSET")
		raven.CaptureMessage("Redis down using inMemory MSET", nil)
		err = c.backup.MSet(pairs, ttl)
		return err
	}

	// MSET doesn't support an expiration, so if an expiration was defined, set it
	// manually. A pipeline is used to send all the commands in one request. The
	// commands are not executed all in one atomic transaction to avoid blocking
	// incoming requests.
	if ttl != 0 {
		pipe := c.client.Pipeline()
		for key := range pairs {
			pipe.Expire(key, ttl)
		}
		_, expErr := pipe.Exec()
		if expErr != nil {
			log.Println("Error setting keys expiration")
			raven.CaptureError(expErr, nil)
		}

	}

	return err
}

func (c *RedisCache) cacheKey(key string) string {
	lowerKey := strings.ToLower(key)

	return fmt.Sprintf("%v-v%v", lowerKey, c.version)
}
