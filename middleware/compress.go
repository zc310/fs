package middleware

import (
	"github.com/valyala/fasthttp"
)

type Compress struct{}

func (p *Compress) Init(c *Config) (err error) { return nil }
func (p *Compress) UnInit()                    {}
func (p *Compress) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.CompressHandler(h)
}
func (p *Compress) Handler() fasthttp.RequestHandler { return func(ctx *fasthttp.RequestCtx) {} }
