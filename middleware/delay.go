package middleware

import (
	"time"

	"github.com/valyala/fasthttp"
)

type Delay struct {
	Delay   string `json:"delay"`
	timeout time.Duration
}

func (p *Delay) Init(c *Config) (err error) {
	if p.Delay != "" {
		p.timeout, err = time.ParseDuration(p.Delay)
		return err
	}
	return nil
}
func (p *Delay) UnInit() {}
func (p *Delay) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {}
}
func (p *Delay) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		delay := string(ctx.QueryArgs().Peek("delay"))
		if delay != "" {
			d, err := time.ParseDuration(delay)
			if err == nil {
				time.Sleep(d)
			}
		} else {
			if p.timeout > 0 {
				time.Sleep(p.timeout)
			}
		}
		h(ctx)
	}
}
