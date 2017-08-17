package middleware

import "github.com/valyala/fasthttp"

type Expires struct{}

func (p *Expires) Init(c *Config) (err error) { return nil }
func (p *Expires) UnInit()                    {}
func (p *Expires) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

	}
}
func (p *Expires) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}
