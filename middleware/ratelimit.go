package middleware

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

type Ratelimit struct {
	Methods  map[string]string `json:"methods"`
	Headers  map[string]string `json:"headers"`
	Args     map[string]string `json:"args"`
	Max      int               `json:"max"`
	Duration string            `json:"duration"`
	Timeout  string            `json:"timeout"`
	Message  struct {
		ContentType string `json:"content_type"`
		Message     string `json:"message"`
		StatusCode  string `json:"status_code"`
	} `json:"message"`
	lim *rate.Limiter
}

func (p *Ratelimit) Init(c *Config) (err error) {
	if p.Max <= 0 {
		return errors.New("error rate max")
	}
	ts, err := time.ParseDuration(p.Duration)
	if err != nil {
		return fmt.Errorf("error rate duration %s", p.Duration)
	}
	p.lim = rate.NewLimiter(rate.Every(ts/time.Duration(p.Max)), p.Max)
	return nil
}
func (p *Ratelimit) UnInit() {}

func (p *Ratelimit) Handler() fasthttp.RequestHandler { return func(ctx *fasthttp.RequestCtx) {} }

func (p *Ratelimit) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ct, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := p.lim.Wait(ct)
		if err != nil {
			ctx.Error("rate", 500)
			return
		}
		h(ctx)
	}
}
