package cache

import (
	"time"
)

//Cache interface
type Cache interface {
	// Get get cached value by key
	Get(key string) (value []byte, ok bool)
	// GetRange
	GetRange(key string, low, high int64) (value []byte, ok bool)
	// Put set cached value with key
	Set(key string, value []byte)
	// Put set cached value with key and expire time
	SetTimeout(key string, value []byte, timeout time.Duration)error
	// Delete delete cached value by key
	Delete(key string)
	// ClearAll clear all cache
	ClearAll() (err error)
}
