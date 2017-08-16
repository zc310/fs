package filecache

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCache_Get(t *testing.T) {
	cachepath := filepath.Join(os.TempDir(), "cache")
	cache, err := New(make(map[string]interface{}))
	k := "a"
	assert.Equal(t, nil, err)
	_, ok := cache.Get("aa")
	assert.Equal(t, ok, false)
	err = cache.Set(k, []byte("a"), time.Second*5)
	assert.Equal(t, nil, err)
	b, ok := cache.Get(k)
	assert.Equal(t, ok, true)
	assert.Equal(t, b, []byte("a"))

	k = "b"
	v := bytes.Repeat([]byte("0123456789"), 1024*1024)
	err = cache.Set(k, v, time.Minute*10)
	assert.Equal(t, err, nil)
	b, ok = cache.Get(k)
	assert.Equal(t, ok, true)
	assert.Equal(t, b, v)

	b, ok = cache.GetRange(k, 6, 9)
	assert.Equal(t, ok, true)
	assert.Equal(t, b, []byte("678"))
	err = cache.Delete(k)
	cache.ClearAll()
	os.RemoveAll(cachepath)
}

func BenchmarkCache_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		cache, err := New(nil)
		assert.Equal(b, err, nil)
		cache.Set("a", []byte("01234567890"), time.Hour*10)
		for pb.Next() {
			cache.Get("a")
		}
		cache.ClearAll()
	})
}
func BenchmarkCache_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		cache, _ := New(nil)
		for pb.Next() {
			cache.Set("a", []byte("01234567890"), time.Hour*10)
		}
		cache.ClearAll()
	})

}
