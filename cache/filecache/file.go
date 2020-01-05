package filecache

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"bytes"

	"github.com/dustin/go-humanize"
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
	dir      string
	db       *leveldb.DB
	fileSize uint64
}

func New(store map[string]interface{}) (cache.Cache, error) {
	var dir string
	if store != nil {
		dir = utils.GetString(store["path"])
	}

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
	var m uint64
	if store != nil {
		m, err = humanize.ParseBytes(utils.GetString(store["size"]))
	}
	if m <= 0 {
		m = 128 * 1024
	}
	var c Cache
	err = c.open(dir, m)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (p *Cache) open(dir string, s uint64) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return &os.PathError{Op: "open", Path: dir, Err: fmt.Errorf("not a directory")}
	}
	for i := 0; i < 256; i++ {
		name := filepath.Join(dir, fmt.Sprintf("%02x", i))
		if err := os.MkdirAll(name, 0777); err != nil {
			return err
		}
	}
	p.db, err = leveldb.OpenFile(filepath.Join(dir, "db"), nil)
	if err != nil {
		return err
	}
	p.dir = dir
	p.fileSize = s
	return nil

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
	if cv.Expired() {
		ok = false
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
	if cv.Expired() {
		ok = false
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
func (p *Cache) Set(key string, value []byte) {
	p.SetTimeout(key, value, time.Hour*24*256)
}

// Put set cached value with key and expire time
func (p *Cache) SetTimeout(key string, value []byte, timeout time.Duration) (err error) {
	var cv cacheValue
	var b []byte

	if uint64(len(value)) <= p.fileSize {
		cv.Value = value
	} else {
		hash := sha1.New()
		hash.Write([]byte(key))
		file := hex.EncodeToString(hash.Sum(nil))
		dir := filepath.Join(p.dir, file[0:2])

		cv.FileName = filepath.Join(dir, file)
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
func (p *Cache) Delete(key string) {
	var cv *cacheValue
	var ok bool
	if cv, ok = p.getValue(key); !ok {
		return
	}
	if len(cv.FileName) > 0 {
		if err := os.Remove(cv.FileName); err != nil {
			return
		}
	}
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
	return p.open(p.dir, p.fileSize)
}
