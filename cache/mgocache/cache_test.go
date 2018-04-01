package mgocache

import (
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMgoCache_Get(t *testing.T) {
	s, err := mgo.Dial("mongodb://127.0.0.1:27017/")
	assert.Equal(t, err, nil)
	c := New(s.DB("cache"), "cache", DefaultKey)
	k := "abc"
	v := []byte("0123456789")
	c.SetTimeout(k, v, time.Hour*24)
	b, ok := c.Get(k)
	assert.Equal(t, ok, true)
	assert.Equal(t, b, v)
	b, ok = c.GetRange(k, 6, 9)
	assert.Equal(t, ok, true)
	assert.Equal(t, b, []byte("678"))

	c.SetTimeout(k, []byte("abc"), time.Hour*24)
	b, ok = c.Get(k)
	assert.Equal(t, ok, true)
	assert.Equal(t, b, []byte("abc"))

	c.Delete(k)
	b, ok = c.Get(k)
	assert.Equal(t, ok, false)

}
