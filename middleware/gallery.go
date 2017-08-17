package middleware

import "github.com/valyala/fasthttp"

type Gallery struct{}

func (p *Gallery) Init(c *Config) (err error) { return nil }
func (p *Gallery) UnInit()                    {}
func (p *Gallery) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

	}
}
func (p *Gallery) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}
