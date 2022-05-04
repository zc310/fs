package middleware

import (
	"fmt"
	"github.com/valyala/fasthttp"

	"net/http"
	"runtime"

	"github.com/zc310/log"
)

type Recover struct {
	Log log.Logger
}

func (p *Recover) Init(c *Config) (err error) {
	p.Log = c.Logger.NewWithPrefix("recover")
	return nil
}
func (p *Recover) UnInit()                          {}
func (p *Recover) Handler() fasthttp.RequestHandler { return func(ctx *fasthttp.RequestCtx) {} }
func (p *Recover) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r := r.(type) {
				case error:
					err = r
				default:
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, 4<<10)
				length := runtime.Stack(stack, true)
				p.Log.Errorf("[%s] %s %s\n", "PANIC RECOVER", err, stack[:length])
				ctx.Error(err.Error(), http.StatusInternalServerError)
			}
		}()
		h(ctx)

	}
}
