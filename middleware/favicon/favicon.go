package favicon

import (
	"github.com/valyala/fasthttp"
	"github.com/zc310/fs/middleware"

	"github.com/zc310/utils/fasthttputil"
	"github.com/zc310/utils/fasthttputil/favicon"
)

type Favicon struct{}

func (p *Favicon) Init(c *middleware.Config) (err error) { return nil }
func (p *Favicon) UnInit()                               {}
func (p *Favicon) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) { h(ctx) }
}
func (p *Favicon) Handler() fasthttp.RequestHandler { return favicon.Handler }

type Ok struct{}

func (p *Ok) Init(c *middleware.Config) (err error) { return nil }
func (p *Ok) UnInit()                               {}

func (p *Ok) Handler() fasthttp.RequestHandler { return fasthttputil.Ok }

func (p *Ok) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) { h(ctx) }
}
