package middleware

import "github.com/zc310/utils/hash"

// GetHashFunc get hash func
func GetHashFunc(a string) (f func(b []byte) string) {
	if a == "" {
		return func(b []byte) string { return string(b) }
	}

	if t, ok := hash.Get(a); ok {
		return func(b []byte) string { return t(b) }
	}
	return func(b []byte) string { return hash.MD5(b) }
}
