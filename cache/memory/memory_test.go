package memory

import (
	"testing"
	"time"
)

func BenchmarkCache_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		cache, _ := New(nil)
		cache.Set("a", []byte("01234567890"), time.Hour*10)
		for pb.Next() {
			cache.Get("a")
		}
	})
}
func BenchmarkCache_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		cache, _ := New(nil)
		for pb.Next() {
			cache.Set("a", []byte("01234567890"), time.Hour*10)
		}
	})
}
