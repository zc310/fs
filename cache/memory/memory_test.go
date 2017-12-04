package memory

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCache_Get(t *testing.T) {
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
}
func BenchmarkCache_Get(b *testing.B) {
	b.StopTimer()
	m, _ := New(nil)
	var key [8]byte
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		m.Set(string(key[:]), make([]byte, 8), time.Hour*10)
	}
	b.StartTimer()
	var hitCount int64
	var ok bool
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		if _, ok = m.Get(string(key[:])); ok {
			hitCount++
		}
	}
}
func BenchmarkCache_Set(b *testing.B) {
	m, _ := New(nil)
	var key [8]byte
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		m.Set(string(key[:]), make([]byte, 8), time.Hour*10)
	}
}

func BenchmarkMapSet(b *testing.B) {
	m := make(map[string][]byte)
	var key [8]byte
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		m[string(key[:])] = make([]byte, 8)
	}
}
func BenchmarkMapGet(b *testing.B) {
	b.StopTimer()
	m := make(map[string][]byte)
	var key [8]byte
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		m[string(key[:])] = make([]byte, 8)
	}
	b.StartTimer()
	var hitCount int64
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		if m[string(key[:])] != nil {
			hitCount++
		}
	}
}
