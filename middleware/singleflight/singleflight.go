package singleflight

import (
	"github.com/golang/groupcache/singleflight"
	"github.com/valyala/fasthttp"

	"github.com/zc310/fs/config"
	"github.com/zc310/fs/middleware"

	"github.com/zc310/log"
	"github.com/zc310/utils"

	"github.com/zc310/fs/template"
)

// Singleflight Singleflight
type Singleflight struct {
	Key     string           `json:"key"`
	Timeout string           `json:"timeout"`
	Hash    string           `json:"hash"`
	Check   config.CheckList `json:"check"`
	hashFun func(b []byte) string
	log     log.Logger
	g       singleflight.Group
}

// Init ...
func (p *Singleflight) Init(c *middleware.Config) (err error) {
	p.log = c.Logger.NewWithPrefix("singleflight")
	p.Key = utils.IfEmpty(p.Key, "{request_uri}")
	p.hashFun = middleware.GetHashFunc(p.Hash)
	return nil
}

// UnInit ...
func (p *Singleflight) UnInit() {}

// Handler ...
func (p *Singleflight) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

	}
}

// Process ...
func (p *Singleflight) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var (
			key  []byte
			body interface{}
			err  error
		)
		type Cache struct {
			header fasthttp.ResponseHeader
			body   []byte
		}
		tpl := template.Get()
		defer template.Put(tpl)

		tpl.SetCtx(ctx)
		if ok, err := config.CheckHit(p.Check, tpl); !ok {
			if err != nil {
				p.log.Error(err, ctx.Request.String())
			}
			h(ctx)
			return
		}

		if key, err = tpl.Execute(p.Key); err != nil {
			p.log.Error(err, ctx.Request.String())
			return
		}
		var cache Cache
		body, err = p.g.Do(string(p.hashFun(key)), func() (interface{}, error) {
			h(ctx)
			cache.body = ctx.Response.Body()
			cache.header = ctx.Response.Header
			return cache, nil
		})
		if err != nil {
			p.log.Error(err, ctx.Request.String())
			h(ctx)
			return
		}

		cache = body.(Cache)
		ctx.Response.Header = cache.header
		ctx.Response.SetBody(cache.body)
	}
}
