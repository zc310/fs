package middleware

import (
	"github.com/valyala/fasthttp"
	"github.com/zc310/fs/cache"
	"github.com/zc310/fs/cache/filecache"
	"github.com/zc310/fs/template"
	"github.com/zc310/log"
	"github.com/zc310/utils"
	"github.com/zc310/utils/hash"

	"bufio"
	"bytes"
	"time"

	"github.com/zc310/fs/cache/memory"
	"github.com/zc310/utils/fasthttputil"
	"strconv"
)

type Cache struct {
	Store   map[string]interface{} `json:"store"`
	Key     string                 `json:"key"`
	Timeout string                 `json:"timeout"`
	Hash    string                 `json:"hash"`
	Check   CheckList              `json:"check"`
	timeout time.Duration
	hashFun func(b []byte) []byte
	cache   cache.Cache

	log log.Logger
}

func (p *Cache) Init(c *Config) (err error) {
	newcache := memory.New
	name := utils.GetString(p.Store["name"])
	if name == "" || name == "file" {
		newcache = filecache.New
	}
	if name == "memory" {
		newcache = memory.New
	}
	if p.cache, err = newcache(p.Store); err != nil {
		return
	}
	if p.timeout, err = time.ParseDuration(p.Timeout); err != nil {
		return err
	}

	p.Key = utils.IfEmpty(p.Key, "{request_uri}")
	p.log = c.Logger.NewWithPrefix("cache")
	p.hashFun = GetHashFunc(p.Hash)
	return nil
}
func (p *Cache) UnInit() {}
func (p *Cache) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var (
			key []byte
			err error
		)
		tpl := template.Get()
		defer template.Put(tpl)
		tpl.SetCtx(ctx)

		if ok, err := CheckHit(p.Check, tpl); !ok {
			if err != nil {
				p.log.Error(err, ctx.Request.String())
			}
			h(ctx)
			return
		}

		if key, err = tpl.Execute(p.Key); err != nil {
			p.log.Error(err, ctx.Request.String())
		} else {

			if b, ok := p.cache.Get(p.hashFun(key)); ok {
				ctx.Response.Read(bufio.NewReader(bytes.NewBuffer(b)))

				tag1 := ctx.Response.Header.Peek("ETag")
				if len(tag1) > 0 && bytes.Equal(ctx.Request.Header.Peek("If-None-Match"), tag1) {
					ctx.NotModified()
				}
				return
			}
		}
		h(ctx)
		if age, ok := fasthttputil.GetResponseAge(ctx, p.timeout); ok {
			ctx.Response.Header.Set("Cache-Control", "public, max-age="+strconv.Itoa( int(p.timeout.Seconds())))
			p.cache.Set(key, []byte(ctx.Response.String()), age)
		}

	}
}
func (p *Cache) Handler() fasthttp.RequestHandler { return func(ctx *fasthttp.RequestCtx) {} }

// GetHashFunc get hash func
func GetHashFunc(a string) (f func(b []byte) []byte) {
	if a == "" {
		return func(b []byte) []byte { return b }
	}

	if t, ok := hash.Get(a); ok {
		return func(b []byte) []byte { return []byte(t(b)) }
	}
	return func(b []byte) []byte { return []byte(hash.MD5(b)) }
}
