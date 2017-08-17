package middleware

import "github.com/valyala/fasthttp"

type Fastcgi struct {
	Root string `json:"root"`
}

func (p *Fastcgi) Init(c *Config) (err error) { return nil }
func (p *Fastcgi) UnInit()                    {}
func (p *Fastcgi) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

	}
}
func (p *Fastcgi) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}
