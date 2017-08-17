package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type Cache struct {
	pool *redis.Pool
}

func (p *Cache) Get(key string) (value []byte, ok bool) {
	return nil, false
}

// GetRange
func (p *Cache) GetRange(key string, low, high int64) (value []byte, ok bool) {
	return nil, false

}

// Put set cached value with key and expire time
func (p *Cache) Set(key string, value []byte, timeout time.Duration) (err error) {
	return
}

// Delete delete cached value by key
func (p *Cache) Delete(key string) (err error) {
	return
}

// ClearAll clear all cache
func (p *Cache) ClearAll() (err error) {
	return
}
