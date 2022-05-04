package middleware

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type Helloworld struct {
	hit  map[string]int64
	lock sync.Mutex
}

func (p *Helloworld) Init(c *Config) (err error) {
	p.hit = map[string]int64{}
	return nil
}
func (p *Helloworld) UnInit() {}
func (p *Helloworld) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Hello, world!\n\n%d\n", p.hit[string(ctx.Request.URI().Path())])
		ctx.SetContentType("text/plain; charset=utf8")
		fmt.Fprintf(ctx, "%s\n\n", &ctx.Request)
		var c fasthttp.Cookie
		c.SetKey("cookie-name")
		c.SetValue(time.Now().String())
		ctx.Response.Header.SetCookie(&c)
		p.lock.Lock()
		p.hit[string(ctx.Request.URI().Path())]++
		b, _ := json.MarshalIndent(p.hit, " ", " ")
		fmt.Fprintf(ctx, "\n\n%s", string(b))
		p.lock.Unlock()

	}
}

func (p *Helloworld) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) { h(ctx) }
}
