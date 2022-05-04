package basicauth_test

import (
	"encoding/base64"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/stretchr/testify/assert"
	. "github.com/zc310/fs/middleware/basicauth"
	"net/http"
)

func TestBasicAuth_Process(t *testing.T) {
	auth := &BasicAuth{Auth: map[string]string{"admin": "password"}}
	ctx := new(fasthttp.RequestCtx)
	assert.Equal(t, auth.Authenticate(ctx), false)
	auth.Process(func(ctx *fasthttp.RequestCtx) {})(ctx)
	assert.Equal(t, ctx.Response.StatusCode(), http.StatusUnauthorized)

	ctx.Request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
	assert.Equal(t, auth.Authenticate(ctx), true)
}
