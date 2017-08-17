package middleware

import "github.com/valyala/fasthttp"

type BasicAuth struct {
}

func (p *BasicAuth) Init(c *Config) (err error) { return nil }
func (p *BasicAuth) UnInit()                    {}
func (p *BasicAuth) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

	}
}
func (p *BasicAuth) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}
