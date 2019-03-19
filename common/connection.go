package common

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"time"
)

// GetNewRedisPool return new Redis pool
func GetNewRedisPool() *redis.Pool {
	host := os.Getenv("REDIS_HOST")
	log.Println("Geting redis pool from ", host)
	redisPool := &redis.Pool{
		MaxIdle:     2,
		IdleTimeout: 60 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host+":6379")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) > time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return redisPool
}
