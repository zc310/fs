package middleware

import (
	"github.com/golang/groupcache/singleflight"
	"github.com/valyala/fasthttp"
	"github.com/zc310/fs/template"
	"github.com/zc310/log"
	"github.com/zc310/utils"
)

// Singleflight Singleflight
type Singleflight struct {
	Key     string `json:"key"`
	Timeout string `json:"timeout"`
	Hash    string `json:"hash"`
	hashFun func(b []byte) string
	log     log.Logger
	g       singleflight.Group
}

// Init ...
func (p *Singleflight) Init(c *Config) (err error) {
	p.log = c.Logger.NewWithPrefix("singleflight")
	p.Key = utils.IfEmpty(p.Key, "{request_uri}")
	p.hashFun = GetHashFunc(p.Hash)
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
			key  string
			body interface{}
			err  error
		)
		tpl := template.Get()
		defer template.Put(tpl)
		tpl.SetCtx(ctx)
		if key, err = tpl.Execute(p.Key); err != nil {
			p.log.Error(err, ctx.Request.String())
			return
		}

		body, err = p.g.Do(p.hashFun([]byte(key)), func() (interface{}, error) {
			h(ctx)
			return ctx.Response.String(), nil
		})
		if err != nil {
			p.log.Error(err, ctx.Request.String())
			h(ctx)
			return
		}
		ctx.Response.SetBodyString(body.(string))
	}
}
