package cache

import (
	"time"
)

//Cache interface
type Cache interface {
	// Get get cached value by key
	Get(key []byte) (value []byte, ok bool)
	// GetRange
	GetRange(key []byte, low, high int64) (value []byte, ok bool)
	// Put set cached value with key and expire time
	Set(key []byte, value []byte, timeout time.Duration) (err error)
	// Delete delete cached value by key
	Delete(key []byte) (err error)
	// ClearAll clear all cache
	ClearAll() (err error)
}
