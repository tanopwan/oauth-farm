package common

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

const prefix = "sss_"

// RedisCache implement Cache interface
type RedisCache struct {
	pool *redis.Pool
}

// NewRedisCache return new instance
func NewRedisCache() *RedisCache {
	return &RedisCache{
		pool: GetNewRedisPool(),
	}
}

// Get value from memory cache
func (c *RedisCache) Get(key string) (interface{}, error) {
	db := c.pool.Get()
	defer db.Close()

	hash := HashSHA256(prefix + key)
	value, err := redis.String(db.Do("GET", hash))
	if err == redis.ErrNil {
		// return empty session
		return nil, err
	}

	return value, nil
}

// Del value from memory cache
func (c *RedisCache) Del(key string) error {
	db := c.pool.Get()
	defer db.Close()

	hash := HashSHA256(prefix + key)
	_, err := redis.String(db.Do("DEL", hash))
	if err == redis.ErrNil {
		// return empty session
		return err
	}

	return nil
}

// Set value to memory cache
func (c *RedisCache) Set(key string, value interface{}) error {
	db := c.pool.Get()
	defer db.Close()

	hash := HashSHA256(prefix + key)
	_, err := db.Do("SET", hash, value)
	if err != nil {
		return err
	}
	return nil
}

// SetExpire value to memory cache with expiration
func (c *RedisCache) SetExpire(key string, value interface{}, duration time.Duration) error {
	db := c.pool.Get()
	defer db.Close()

	hash := HashSHA256(prefix + key)
	_, err := db.Do("SETEX", hash, duration, value)
	if err != nil {
		return err
	}
	return nil
}
