package leveldb

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/zc310/fs/cache"
	"github.com/zc310/utils"
)

type cacheValue struct {
	Value   []byte
	Expires time.Time
}

func (p *cacheValue) Expired() bool {
	return p.Expires.Before(time.Now())
}

type Cache struct {
	dir string
	db  *leveldb.DB
}

func New(store map[string]interface{}) (cache.Cache, error) {
	var dir string
	if store != nil {
		dir = utils.GetString(store["path"])
	}

	return NewFromPath(dir)
}
func NewFromPath(dir string) (cache.Cache, error) {

	var err error
	if dir == "" {
		dir, err = ioutil.TempDir(os.TempDir(), "cache")
		if err != nil {
			return nil, err
		}
	} else {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
	}

	var c Cache
	err = c.open(dir)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (p *Cache) open(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return &os.PathError{Op: "open", Path: dir, Err: fmt.Errorf("not a directory")}
	}

	p.db, err = leveldb.OpenFile(dir, nil)
	if err != nil {
		return err
	}
	p.dir = dir

	return nil

}
func (p *Cache) Close() error {
	return p.db.Close()
}
func (p *Cache) getValue(key string) (*cacheValue, bool) {
	b, err := p.db.Get([]byte(key), nil)
	if err != nil {
		return nil, false
	}

	r := new(cacheValue)
	if err = msgpack.Unmarshal(b, r); err != nil {
		return nil, false
	}
	return r, true
}

// Get get cached value by key
func (p *Cache) Get(key string) (value []byte, ok bool) {
	var cv *cacheValue
	if cv, ok = p.getValue(key); !ok {
		return
	}
	if cv.Expired() {
		ok = false
		return
	}

	value = cv.Value

	return
}

// GetRange
func (p *Cache) GetRange(key string, low, high int64) (value []byte, ok bool) {
	if high == 0 {
		return p.Get(key)
	}

	var cv *cacheValue
	if cv, ok = p.getValue(key); !ok {
		return
	}
	if cv.Expired() {
		ok = false
		return
	}

	value = cv.Value[low:high]

	return
}
func (p *Cache) Set(key string, value []byte) {
	p.SetTimeout(key, value, time.Hour*24*256)
}

// Put set cached value with key and expire time
func (p *Cache) SetTimeout(key string, value []byte, timeout time.Duration) (err error) {
	var cv cacheValue
	var b []byte

	cv.Value = value

	cv.Expires = time.Now().Add(timeout)
	if b, err = msgpack.Marshal(cv); err != nil {
		return err
	}

	return p.db.Put([]byte(key), b, nil)
}

// Delete delete cached value by key
func (p *Cache) Delete(key string) {
	p.db.Delete([]byte(key), nil)

}

// ClearAll clear all cache
func (p *Cache) ClearAll() (err error) {
	err = p.db.Close()
	if err != nil {
		return err
	}
	err = os.RemoveAll(p.dir)
	if err != nil {
		return err
	}
	return p.open(p.dir)
}
