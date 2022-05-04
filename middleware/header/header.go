package header

import (
	"github.com/valyala/fasthttp"
	"github.com/zc310/fs/middleware"
)

type Header struct {
	Add map[string]string `json:"add"`
	Set map[string]string `json:"set"`
	Del []string          `json:"del"`
}

func (p *Header) Init(c *middleware.Config) (err error) { return nil }
func (p *Header) UnInit()                               {}
func (p *Header) Handler() fasthttp.RequestHandler      { return func(ctx *fasthttp.RequestCtx) {} }
func (p *Header) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
		for k, v := range p.Set {
			ctx.Response.Header.Set(k, v)
		}
		for k, v := range p.Add {
			ctx.Response.Header.Add(k, v)
		}
		for _, k := range p.Del {
			ctx.Response.Header.Del(k)
		}
	}
}
