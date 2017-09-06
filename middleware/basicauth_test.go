package middleware

import (
	"encoding/base64"
	"testing"

	"github.com/valyala/fasthttp"

	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth_Process(t *testing.T) {
	auth := &BasicAuth{Auth: map[string]string{"admin": "password"}}
	ctx := new(fasthttp.RequestCtx)
	assert.Equal(t, auth.authenticate(ctx), false)
	auth.Process(func(ctx *fasthttp.RequestCtx) {})(ctx)
	assert.Equal(t, ctx.Response.StatusCode(), http.StatusUnauthorized)

	ctx.Request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	assert.Equal(t, auth.authenticate(ctx), true)
}
