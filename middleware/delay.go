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
		var err error
		delay := string(ctx.QueryArgs().Peek("delay"))
		timeout := p.timeout
		if delay != "" {
			timeout, err = time.ParseDuration(delay)
			if err != nil {
				h(ctx)
			}
		} else {
			if p.timeout <= 0 {
				h(ctx)
			}
		}

		for {
			select {
			case <-time.After(timeout):
				h(ctx)
				return
			}
		}

	}
}
