package storage

import (
	"log"
	"net"
	"regexp"
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
	client *redisTrace.Client
	backup InMemoryCache
}

// ErrNotFound is returned when an item is not present in cache
var ErrNotFound = errors.New("item not found")

var redisDownRegex = regexp.MustCompile(`^dial tcp .* connect: connection refused|dial tcp: lookup .* no such host|EOF$`)

// isConnectionRefused returns whether the error passed as an argument is a redis
// connection error or not.
func isConnectionRefused(err error) bool {
	return err != nil && redisDownRegex.Match(([]byte(err.Error())))
}

// NewRedisCache initializes and returns a RedisCache
func NewRedisCache(host, port string) *RedisCache {
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
	lowerKey := strings.ToLower(key)
	data, err := c.client.Get(lowerKey).Bytes()
	if err == nil {
		err = json.Unmarshal(data, &value)
		return err
	}
	if err == redis.Nil {
		return ErrNotFound
	}

	if isConnectionRefused(err) {
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

	lowerKey := strings.ToLower(key)
	_, err = c.client.Set(lowerKey, strVal, ttl).Result()
	if isConnectionRefused(err) {
		log.Println("Redis down using inMemory SET")
		raven.CaptureMessage("Redis down using inMemory SET", nil)
		err = c.backup.Set(key, value)
	}
	return err
}

// Del deletes an item from cache
func (c *RedisCache) Del(key string) error {
	lowerKey := strings.ToLower(key)
	_, err := c.client.Del(lowerKey).Result()
	if isConnectionRefused(err) {
		log.Println("Redis down using inMemory DEL")
		raven.CaptureMessage("Redis down using inMemory DEL", nil)
		c.backup.Del(key)
		return nil
	}
	return err
}

// MSet writes items to cache in bulk
func (c *RedisCache) MSet(pairs map[string]interface{}) error {
	args := make([]interface{}, len(pairs)*2)
	i := 0

	for key, value := range pairs {
		args[i] = strings.ToLower(key)
		strValue, err := json.Marshal(value)
		if err != nil {
			return err
		}
		args[i+1] = strValue
		i += 2
	}

	_, err := c.client.MSet(args...).Result()
	if isConnectionRefused(err) {
		log.Println("Redis down using inMemory MSET")
		raven.CaptureMessage("Redis down using inMemory MSET", nil)
		err = c.backup.MSet(pairs)
	}
	return err
}