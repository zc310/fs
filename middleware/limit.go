package middleware

import (
	"github.com/valyala/fasthttp"
)

type Limit struct{}

func (p *Limit) Init(c *Config) (err error) { return nil }
func (p *Limit) UnInit()                    {}
func (p *Limit) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

	}
}
func (p *Limit) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}
