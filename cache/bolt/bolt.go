package bolt

import (
	"github.com/vmihailenco/msgpack/v5"
	"github.com/zc310/fs/cache"
	"github.com/zc310/utils"
	"go.etcd.io/bbolt"
	"io/ioutil"
	"os"
	"time"
)

type cacheValue struct {
	Value   []byte
	Expires time.Time
}

func (p *cacheValue) Expired() bool {
	return p.Expires.Before(time.Now())
}

var bucket = []byte("a")

type Cache struct {
	dir string
	db  *bbolt.DB
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
	}

	var c Cache
	err = c.open(dir)
	if err != nil {
		return nil, err
	}

	return &c, c.db.Update(func(tx *bbolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(bucket)
		return err
	})
}

func (p *Cache) open(file string) error {

	var err error
	p.db, err = bbolt.Open(file, 0666, nil)
	if err != nil {
		return err
	}
	p.dir = file

	return nil

}
func (p *Cache) Close() error {
	return p.db.Close()
}
func (p *Cache) getValue(key string) (*cacheValue, bool) {
	var b []byte
	err := p.db.View(func(tx *bbolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return bbolt.ErrBucketNotFound
		}
		b = bk.Get([]byte(key))
		return nil
	})

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

	return p.db.Update(func(tx *bbolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return bbolt.ErrBucketNotFound
		}

		return bk.Put([]byte(key), b)
	})
}

// Delete delete cached value by key
func (p *Cache) Delete(key string) {
	_ = p.db.Update(func(tx *bbolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return bbolt.ErrBucketNotFound
		}

		return bk.Delete([]byte(key))
	})

}

// ClearAll clear all cache
func (p *Cache) ClearAll() (err error) {
	err = p.db.Close()
	if err != nil {
		return err
	}
	err = os.Remove(p.dir)
	if err != nil {
		return err
	}
	return p.open(p.dir)
}
