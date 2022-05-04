package middleware_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/zc310/fs/middleware"
	"testing"

	"github.com/zc310/log"
	"net/http"
)

func TestRecoverHandler(t *testing.T) {
	var ctx fasthttp.RequestCtx

	h := func(ctx *fasthttp.RequestCtx) {
		panic("test")
	}

	r := new(middleware.Recover)
	c := new(middleware.Config)
	c.Logger = log.NewWithPrefix("recover")
	r.Init(c)
	r.Process(h)(&ctx)
	assert.Equal(t, http.StatusInternalServerError, ctx.Response.StatusCode())

}
