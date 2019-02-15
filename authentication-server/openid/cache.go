package openid

import (
	"fmt"
	"time"
)

type cache interface {
	get(key string) interface{}
	del(key string)
	set(key string, value interface{})
	setExpire(key string, value interface{}, duration time.Duration)
}

type memoryCache struct {
	data map[string]interface{}
}

func newMemoryCache() *memoryCache {
	return &memoryCache{
		data: make(map[string]interface{}),
	}
}

func (c *memoryCache) get(key string) interface{} {
	return c.data[key]
}

func (c *memoryCache) del(key string) {
	delete(c.data, key)
}

func (c *memoryCache) set(key string, value interface{}) {
	c.data[key] = value
}

func (c *memoryCache) setExpire(key string, value interface{}, duration time.Duration) {
	c.data[key] = value
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		delete(c.data, key)
		fmt.Println("google_public cache is expired")
	}()
}
