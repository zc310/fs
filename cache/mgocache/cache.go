package mgocache

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"io"
	"time"
)

type cacheInfo struct {
	Expire time.Time `bson:"e"`
}

type MgoCache struct {
	db  *mgo.Database
	gfs *mgo.GridFS

	fun func(string) string
}

func (p *MgoCache) Get(key string) ([]byte, bool) {
	f, err := p.gfs.OpenId(p.fun(key))
	if err != nil {
		return []byte{}, false
	}
	b := make([]byte, f.Size())
	_, err = f.Read(b)
	var info cacheInfo
	err = f.GetMeta(&info)
	if err != nil {
		return []byte{}, false
	}
	return b, info.Expire.After(time.Now())
}
func (p *MgoCache) ClearAll() (err error) {
	_, err = p.gfs.Files.RemoveAll(bson.M{})
	if err != nil {
		return
	}
	_, err = p.gfs.Chunks.RemoveAll(bson.M{})
	return
}
func (p *MgoCache) GetRange(key string, low, high int64) (b []byte, ok bool) {
	f, err := p.gfs.OpenId(p.fun(key))
	if err != nil {
		return
	}
	_, err = f.Seek(low, io.SeekStart)
	if err != nil {
		return
	}
	b = make([]byte, high-low)
	_, err = f.Read(b)
	if err != nil {
		return []byte{}, false
	}
	var info cacheInfo
	err = f.GetMeta(&info)
	if err != nil {
		return []byte{}, false
	}
	return b, info.Expire.After(time.Now())
}
func (p *MgoCache) Set(key string, value []byte) {
	p.SetTimeout(key, value, time.Hour*24*256)
}
func (p *MgoCache) SetTimeout(key string, value []byte, timeout time.Duration) {
	k := p.fun(key)
	err := p.gfs.Remove(k)
	if err != nil {
		return
	}
	f, err := p.gfs.Create(k)
	if err != nil {
		return
	}

	f.SetMeta(cacheInfo{time.Now().Add(timeout)})
	f.SetId(f.Name())

	_, err = f.Write(value)
	if err != nil {
		return
	}
	f.Close()

}

func (p *MgoCache) Delete(key string) {
	p.gfs.RemoveId(p.fun(key))
}

func Md5key(s string) string {
	h := md5.New()
	io.WriteString(h, s)

	return hex.EncodeToString(h.Sum(nil))
}
func DefaultKey(s string) string {
	return s
}

func New(db *mgo.Database, name string, f func(string) string) *MgoCache {
	return &MgoCache{db, db.GridFS(name), f}
}
