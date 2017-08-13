package filecache

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"time"
	"path/filepath"
	"bytes"
)

func TestCache_Get(t *testing.T) {
	cachepath:=filepath.Join( os.TempDir(),"cache")
	cache,err:=New(make(map[string]interface{}))
	k:="a"
	assert.Equal(t,nil,err)
	_,ok:=cache.Get("aa")
	assert.Equal(t,ok,false)
	err=cache.Set(k,[]byte("a"),time.Second*5)
	assert.Equal(t,nil,err)
	b,ok:=cache.Get(k)
	assert.Equal(t,ok,true)
	assert.Equal(t,b,[]byte("a"))

	k="b"
	v:=bytes.Repeat([]byte("0123456789"),1024*1024)
	err=cache.Set(k,v,time.Minute*10)
	assert.Equal(t,err,nil)
	b,ok=cache.Get(k)
	assert.Equal(t,ok,true)
	assert.Equal(t,b,v)

	b,ok=cache.GetRange(k,6,9)
	assert.Equal(t,ok,true)
	assert.Equal(t,b,[]byte("678"))
	err=cache.Delete(k)
	cache.ClearAll()
	os.RemoveAll(cachepath)
}