package middleware

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

type Helloworld struct {
	i int64
}

func (p *Helloworld) Init(c *Config) (err error) { return nil }
func (p *Helloworld) UnInit()                    {}
func (p *Helloworld) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Hello, world!\n\n%d\n", p.i)
		ctx.SetContentType("text/plain; charset=utf8")
		fmt.Fprintf(ctx, "%s\n\n", &ctx.Request)
		var c fasthttp.Cookie
		c.SetKey("cookie-name")
		c.SetValue(time.Now().String())
		ctx.Response.Header.SetCookie(&c)
		atomic.AddInt64(&p.i, 1)
	}
}

func (p *Helloworld) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) { h(ctx) }
}
