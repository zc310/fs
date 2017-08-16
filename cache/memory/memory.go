package memory

import (
	"github.com/zc310/apiproxy/cache"
	"gopkg.in/vmihailenco/msgpack.v2"
	"sync"
	"time"
)

type cacheValue struct {
	Value   []byte
	Expires time.Time
}

func (p *cacheValue) Expired() bool {
	return p.Expires.Before(time.Now())
}

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
	r := new(cacheValue)
	if err := msgpack.Unmarshal(b, r); err != nil {
		return nil, false
	}
	return r.Value, !r.Expired()
}

// GetRange
func (p *MemoryCache) GetRange(key string, low, high int64) (value []byte, ok bool) {
	if value, ok = p.Get(key); !ok {
		return
	}
	return value[low:high], ok
}

// Put set cached value with key and expire time
func (p *MemoryCache) Set(key string, value []byte, timeout time.Duration) (err error) {
	var b []byte
	b, err = msgpack.Marshal(&cacheValue{value, time.Now().Add(timeout)})
	p.mu.RLock()
	p.items[key] = b
	p.mu.RUnlock()
	return err
}

// Delete delete cached value by key
func (p *MemoryCache) Delete(key string) (err error) {
	p.mu.Lock()
	delete(p.items, key)
	p.mu.Unlock()
	return nil
}

// ClearAll clear all cache
func (p *MemoryCache) ClearAll() (err error) {
	p.mu.Lock()
	p.items = nil
	p.mu.Unlock()
	return nil
}
