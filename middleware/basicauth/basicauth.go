package basicauth

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/zc310/fs/middleware"

	"github.com/zc310/headers"
)

type BasicAuth struct {
	Auth map[string]string
}

func (p *BasicAuth) Init(c *middleware.Config) (err error) { return nil }
func (p *BasicAuth) UnInit()                               {}
func (p *BasicAuth) Handler() fasthttp.RequestHandler      { return func(ctx *fasthttp.RequestCtx) {} }
func (p *BasicAuth) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !p.Authenticate(ctx) {
			ctx.Response.Header.Set("WWW-Authenticate", `Basic realm=Restricted`)
			ctx.SetStatusCode(http.StatusUnauthorized)
			return
		}
		h(ctx)
	}
}
func (p *BasicAuth) Authenticate(ctx *fasthttp.RequestCtx) bool {
	const basicScheme = "Basic "
	auth := string(ctx.Request.Header.Peek(headers.Authorization))
	if !strings.HasPrefix(auth, basicScheme) {
		return false
	}
	str, err := base64.StdEncoding.DecodeString(auth[len(basicScheme):])
	if err != nil {
		return false
	}
	creds := bytes.SplitN(str, []byte(":"), 2)
	if len(creds) != 2 {
		return false
	}

	if pass, ok := p.Auth[string(creds[0])]; ok {
		return pass == string(creds[1])
	}
	return false

}
