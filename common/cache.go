package common

import (
	"time"
)

// Cache interface
type Cache interface {
	Get(key string) (interface{}, error)
	Del(key string) error
	Set(key string, value interface{}) error
	SetExpire(key string, value interface{}, duration time.Duration) error
}
