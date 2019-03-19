package openid

type memoryCache struct {
	data map[string]interface{}
}

func newMemoryCache() *memoryCache {
	return &memoryCache{
		data: make(map[string]interface{}),
	}
}

// Get value from memory cache
func (c *memoryCache) Get(key string) interface{} {
	return c.data[key]
}

// Del value from memory cache
func (c *memoryCache) Del(key string) {
	delete(c.data, key)
}

// Set value to memory cache
func (c *memoryCache) Set(key string, value interface{}) {
	c.data[key] = value
}

// Set value to memory cache with expiration
func (c *memoryCache) SetExpire(key string, value interface{}, duration time.Duration) {
	c.data[key] = value
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		delete(c.data, key)
		fmt.Printf("%s cache is expired", key)
	}()
}
