package middleware

import (
	"github.com/valyala/fasthttp"
	"github.com/zc310/apiproxy/cache"
	"github.com/zc310/apiproxy/cache/filecache"
	"github.com/zc310/apiproxy/template"
	"github.com/zc310/log"
	"github.com/zc310/utils"
	"github.com/zc310/utils/hash"

	"bufio"
	"bytes"
	"time"

	"github.com/zc310/apiproxy/cache/memory"
	"github.com/zc310/utils/fasthttputil"
)

type Cache struct {
	Store   map[string]interface{} `json:"store"`
	Key     string                 `json:"key"`
	Timeout string                 `json:"timeout"`
	Hash    string                 `json:"hash"`
	timeout time.Duration
	hashFun func(b []byte) string
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

	if p.Key == "" {
		p.Key = "{request_uri}"
	}

	p.log = c.Logger.NewWithPrefix("cache")

	var ok bool
	if p.hashFun, ok = hash.Get(p.Hash); !ok {
		p.hashFun = hash.MD5
	}
	return nil
}
func (p *Cache) UnInit() {}
func (p *Cache) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var (
			key string
			err error
		)
		tpl := template.Get()
		defer template.Put(tpl)
		tpl.SetCtx(ctx)
		if key, err = tpl.Execute(p.Key); err != nil {
			p.log.Error(err, ctx.Request.String())
		} else {
			key = p.hashFun([]byte(key))
			if b, ok := p.cache.Get(key); ok {
				ctx.Response.Read(bufio.NewReader(bytes.NewBuffer(b)))
				return

				//
				//etag := ctx.Response.Header.Peek("etag")
				//if len(etag) > 0 && len(ctx.Request.Header.Peek("etag")) == 0 {
				//	lastModified := ctx.Response.Header.Peek("last-modified")
				//	if len(lastModified) > 0 && len(ctx.Request.Header.Peek("last-modified")) == 0 {
				//		ctx.Request.Header.SetBytesV("if-none-match", etag)
				//		ctx.Request.Header.SetBytesV("if-modified-since", lastModified)
				//	}
				//}

			}
		}
		h(ctx)
		p.cache.Set(key, []byte(ctx.Response.String()), fasthttputil.GetResponseAge(ctx, p.timeout))
	}
}
func (p *Cache) Handler() fasthttp.RequestHandler { return func(ctx *fasthttp.RequestCtx) {} }
