package filecache

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"bytes"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zc310/fs/cache"
	"github.com/zc310/utils"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type cacheValue struct {
	Value    []byte
	FileName string
	Expires  time.Time
}

func (p *cacheValue) Expired() bool {
	return p.Expires.Before(time.Now())
}

type Cache struct {
	cachePath string
	db        *leveldb.DB
}

func New(store map[string]interface{}) (cache.Cache, error) {
	var cachepath string
	if store != nil {
		cachepath = utils.GetString(store["path"])
	}
	var err error
	if cachepath == "" {
		cachepath, err = ioutil.TempDir(os.TempDir(), "cache")
		if err != nil {
			return nil, err
		}
	}
	db, err := leveldb.OpenFile(filepath.Join(cachepath, "db"), nil)
	if err != nil {
		return nil, err
	}
	return &Cache{cachepath, db}, nil
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
	var err error
	var cv *cacheValue
	if cv, ok = p.getValue(key); !ok {
		return
	}
	if ok = !cv.Expired(); !ok {
		return
	}

	if len(cv.FileName) != 0 {
		value, err = ioutil.ReadFile(cv.FileName)
		ok = err == nil
	} else {
		value = cv.Value
	}
	return
}

// GetRange
func (p *Cache) GetRange(key string, low, high int64) (value []byte, ok bool) {
	if high == 0 {
		return p.Get(key)
	}
	var err error
	var cv *cacheValue
	if cv, ok = p.getValue(key); !ok {
		return
	}
	if ok = !cv.Expired(); !ok {
		return
	}

	if len(cv.FileName) != 0 {
		var f *os.File
		f, err = os.Open(cv.FileName)
		if err != nil {
			return
		}
		defer f.Close()
		_, err = f.Seek(low, 0)
		buf := new(bytes.Buffer)
		_, err = io.CopyN(buf, f, high-low)
		ok = err == nil
		value = buf.Bytes()
	} else {
		value = cv.Value[low:high]
	}
	return
}

// Put set cached value with key and expire time
func (p *Cache) Set(key string, value []byte, timeout time.Duration) (err error) {
	var cv cacheValue
	var b []byte

	if len(value) <= 10*1024 {
		cv.Value = value
	} else {
		hash := sha1.New()
		io.WriteString(hash, key)
		file := hex.EncodeToString(hash.Sum(nil))
		cachepath := filepath.Join(p.cachePath, filepath.Join(file[0:2], file[7:9]))
		if err = os.MkdirAll(cachepath, os.ModePerm); err != nil {
			return
		}
		cv.FileName = filepath.Join(cachepath, file)
		if err = ioutil.WriteFile(cv.FileName, value, os.ModePerm); err != nil {
			return err
		}
	}
	cv.Expires = time.Now().Add(timeout)
	if b, err = msgpack.Marshal(cv); err != nil {
		return err
	}

	return p.db.Put([]byte(key), b, nil)
}

// Delete delete cached value by key
func (p *Cache) Delete(key string) (err error) {
	var cv *cacheValue
	var ok bool
	if cv, ok = p.getValue(key); !ok {
		return
	}
	if len(cv.FileName) > 0 {
		if err = os.Remove(cv.FileName); err != nil {
			return
		}
	}
	return p.db.Delete([]byte(key), nil)
}

// ClearAll clear all cache
func (p *Cache) ClearAll() (err error) {
	p.db.Close()
	dbpath := filepath.Join(p.cachePath, "db")
	os.RemoveAll(dbpath)
	p.db, err = leveldb.OpenFile(dbpath, nil)
	if err != nil {
		return err
	}
	return nil
}
