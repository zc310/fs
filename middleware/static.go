package middleware

import (
	"github.com/valyala/fasthttp"
	"github.com/zc310/fs/template"
	"strings"
)

type Static struct {
	h                  fasthttp.RequestHandler
	Root               string   `json:"root"`
	Index              []string `json:"index"`
	GenerateIndexPages bool     `json:"generate_index_pages"`
	AcceptByteRange    bool     `json:"byte_range"`
}

func (p *Static) Init(c *Config) (err error) {
	if len(p.Index) == 0 {
		p.Index = []string{"index.html"}
	}
	root, err := template.NewTemplate(p.Root)
	if err != nil {
		return err
	}
	fs := &fasthttp.FS{
		Root:               root,
		IndexNames:         p.Index,
		GenerateIndexPages: p.GenerateIndexPages,
		AcceptByteRange:    p.AcceptByteRange,
	}
	fs.PathRewrite = fasthttp.NewPathSlashesStripper(strings.Count(c.Path, "/") - 1)
	p.h = fs.NewRequestHandler()
	return nil
}
func (p *Static) UnInit() {}
func (p *Static) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}
func (p *Static) Handler() fasthttp.RequestHandler {
	return p.h
}
