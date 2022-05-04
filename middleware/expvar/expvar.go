package expvar

import (
	"expvar"
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/zc310/fs/middleware"
)

type ExpVar struct{}

func (p *ExpVar) Init(c *middleware.Config) (err error) { return nil }
func (p *ExpVar) UnInit()                               {}
func (p *ExpVar) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(ctx, "{\n")
		first := true
		expvar.Do(func(kv expvar.KeyValue) {
			if !first {
				fmt.Fprintf(ctx, ",\n")
			}
			first = false
			fmt.Fprintf(ctx, "%q: %s", kv.Key, kv.Value)
		})
		fmt.Fprintf(ctx, "\n}\n")
	}
}
func (p *ExpVar) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}
