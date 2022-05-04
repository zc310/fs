package gzip

import (
	"github.com/valyala/fasthttp"
	"github.com/zc310/fs/middleware"
)

type Compress struct{}

func (p *Compress) Init(c *middleware.Config) (err error) { return nil }
func (p *Compress) UnInit()                               {}
func (p *Compress) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.CompressHandlerBrotliLevel(h, fasthttp.CompressBrotliDefaultCompression, fasthttp.CompressDefaultCompression)
}
func (p *Compress) Handler() fasthttp.RequestHandler { return func(ctx *fasthttp.RequestCtx) {} }
