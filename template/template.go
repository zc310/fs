package template

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"sync"

	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasttemplate"
	"github.com/zc310/utils/fasthttputil"
)

type Template struct {
	ctx *fasthttp.RequestCtx
	buf bytebufferpool.Pool
}

var pool = &sync.Pool{
	New: func() interface{} {
		return New(nil)
	},
}

func Get() *Template {
	return pool.Get().(*Template)
}
func Put(t *Template) {
	pool.Put(t)
}
func New(ctx *fasthttp.RequestCtx) *Template {
	return &Template{ctx: ctx}
}
func (p *Template) SetCtx(ctx *fasthttp.RequestCtx) {
	p.ctx = ctx
}
func (p *Template) Execute(t string) ([]byte, error) {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	_, err := fasttemplate.New(t, "{", "}").ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "tempdir":
			return w.Write([]byte(os.TempDir()))
		default:
			if p.ctx != nil {
				switch tag {
				case "scheme":
					return w.Write(p.ctx.URI().Scheme())
				case "host":
					return w.Write(p.ctx.Host())
				case "request_uri":
					return w.Write(p.ctx.URI().FullURI())
				case "request_method":
					return w.Write(p.ctx.Method())
				case "request_body":
					return w.Write(p.ctx.Request.Body())
				case "content_type":
					return w.Write(p.ctx.Request.Header.Peek("Content-Type"))
				case "accept_encoding":
					return w.Write(p.ctx.Request.Header.Peek("Accept-Encoding"))
				case "content_length":
					return w.Write([]byte(strconv.Itoa(len(p.ctx.PostBody()))))
				case "query_string":
					return w.Write(p.ctx.Request.URI().QueryString())
				case "document_root":
					return w.Write(p.ctx.Path())
				case "remote_addr":
					return w.Write([]byte(p.ctx.RemoteIP().String()))
				default:
					return w.Write(fasthttputil.GetArgs(p.ctx, tag))
				}
			}
			return 0, fmt.Errorf("[unknown tag %q]", tag)
		}
	})
	return buf.Bytes(), err
}

func NewTemplate(t string) ([]byte, error) {
	return New(nil).Execute(t)
}
