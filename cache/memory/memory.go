package memory

import (
	"encoding/binary"
	"sync"
	"time"

	"github.com/zc310/fs/cache"
)

// New returns a new Cache in in-memory
func New(_ map[string]interface{}) (cache.Cache, error) {
	return &MemoryCache{items: map[string][]byte{}}, nil
}

// MemoryCache  Cache in in-memory map.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string][]byte
}

// Get get cached value by key
func (p *MemoryCache) Get(key string) (value []byte, ok bool) {
	var b []byte
	p.mu.RLock()
	b, ok = p.items[key]
	p.mu.RUnlock()
	if !ok {
		return
	}
	value = make([]byte, len(b)-8)
	copy(value, b[8:])
	return b[8:], binary.LittleEndian.Uint64(b[:8]) > uint64(time.Now().Unix())
}

// GetRange
func (p *MemoryCache) GetRange(key string, low, high int64) (value []byte, ok bool) {
	if value, ok = p.Get(key); !ok {
		return
	}
	return value[low:high], ok
}
func (p *MemoryCache) Set(key string, value []byte) {
	p.SetTimeout(key, value, time.Hour*24*256)
}

// Put set cached value with key and expire time
func (p *MemoryCache) SetTimeout(key string, value []byte, timeout time.Duration) (err error) {
	b := make([]byte, 8+len(value))
	binary.LittleEndian.PutUint64(b, uint64(time.Now().Add(timeout).Unix()))

	copy(b[8:], value)

	p.mu.RLock()
	p.items[string(key)] = b
	p.mu.RUnlock()
	return err
}

// Delete delete cached value by key
func (p *MemoryCache) Delete(key string) {
	p.mu.Lock()
	delete(p.items, key)
	p.mu.Unlock()

}

// ClearAll clear all cache
func (p *MemoryCache) ClearAll() (err error) {
	p.mu.Lock()
	p.items = nil
	p.mu.Unlock()
	return nil
}
func (p *MemoryCache) Close() error {
	return nil
}
