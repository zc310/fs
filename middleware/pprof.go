package middleware

import (
	"github.com/valyala/fasthttp"
	"github.com/zc310/utils/fasthttputil/pprof"
)

type Pprof struct{}

func (p *Pprof) Init(c *Config) (err error)       { return nil }
func (p *Pprof) UnInit()                          {}
func (p *Pprof) Handler() fasthttp.RequestHandler { return pprof.Handler }
func (p *Pprof) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) { h(ctx) }
}
