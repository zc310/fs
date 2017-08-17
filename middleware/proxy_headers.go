package middleware

import (
	"github.com/valyala/fasthttp"
	"github.com/zc310/utils/fasthttputil"
)

type ProxyHeader struct {
	Set map[string]string `json:"set"`
	Add map[string]string `json:"add"`
	Del []string          `json:"del"`
}

func (p *ProxyHeader) Init(c *Config) (err error)       { return nil }
func (p *ProxyHeader) UnInit()                          {}
func (p *ProxyHeader) Handler() fasthttp.RequestHandler { return func(ctx *fasthttp.RequestCtx) {} }
func (p *ProxyHeader) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		for k, v := range p.Set {
			ctx.Request.Header.Set(k, v)
		}
		for k, v := range p.Add {
			ctx.Request.Header.Add(k, v)
		}
		for _, k := range p.Del {
			ctx.Request.Header.Del(k)
		}
		setIP(ctx)
		h(ctx)
	}
}
func setIP(ctx *fasthttp.RequestCtx) {
	r := &ctx.Request
	var s string
	xff := r.Header.Peek(fasthttputil.XForwardedFor)
	xr := r.Header.Peek(fasthttputil.XRealIP)
	if len(xff) == 0 {
		s = ctx.RemoteIP().String()
		if len(xr) == 0 {
			r.Header.Set(fasthttputil.XRealIP, s)
		}
	} else {
		s = string(xff) + ", " + ctx.RemoteIP().String()
	}
	r.Header.Set(fasthttputil.XForwardedFor, s)

	if len(r.Header.Peek(fasthttputil.XForwardedHost)) > 0 {
		r.SetHostBytes(r.Header.Peek(fasthttputil.XForwardedHost))
	}
}
