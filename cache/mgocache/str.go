package mgocache

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/zc310/fs/cache"
)

type Cache struct {
	gfs *mgo.Collection

	fun func(string) string
}
type cacheValue struct {
	Value  []byte    `bson:"v"`
	Expire time.Time `bson:"t"`
}

func (p *Cache) Get(key string) ([]byte, bool) {
	var cv cacheValue
	err := p.gfs.FindId(p.fun(key)).One(&cv)
	if err != nil {
		return []byte{}, false
	}

	return cv.Value, cv.Expire.After(time.Now())
}
func (p *Cache) ClearAll() (err error) {
	_, err = p.gfs.RemoveAll(bson.M{})

	return
}
func (p *Cache) Delete(key string) {
	_ = p.gfs.RemoveId(p.fun(key))

}
func (p *Cache) GetRange(key string, low, high int64) (b []byte, ok bool) {
	if high == 0 {
		return p.Get(key)
	}

	if b, ok = p.Get(key); !ok {
		return
	}
	if !ok {
		return
	}

	b = b[low:high]

	return
}
func (p *Cache) Set(key string, value []byte) {
	p.SetTimeout(key, value, time.Hour*24*256)
}
func (p *Cache) SetTimeout(key string, value []byte, timeout time.Duration) error {
	k := p.fun(key)
	_, err := p.gfs.UpsertId(k, bson.M{"$set": bson.M{"v": value, "t": time.Now().Add(timeout)}})
	return err
}
func (p *Cache) Close() error {
	return nil
}
func NewCache(col *mgo.Collection, f func(string) string) cache.Cache {
	if f == nil {
		f = sha1key
	}
	return &Cache{col, f}
}
