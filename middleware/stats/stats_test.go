package stats_test

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/zc310/fs/middleware"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/zc310/fasthttprouter"
)

func TestStats_Handler(t *testing.T) {
	var ctx fasthttp.RequestCtx

	s := "GET /a HTTP/1.1\nHost: aaa.com\n\n"
	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := ctx.Request.Read(br); err != nil {
		t.Fatalf("cannot read request: %s", err)
	}
	h := func(ctx *fasthttp.RequestCtx) {
		fmt.Fprint(ctx, "ok\n")

	}
	st := new(Stats)
	c := new(Config)
	c.Router = fasthttprouter.New()
	st.Init(c)
	st.Process(h)(&ctx)
	assert.Equal(t, int64(1), st.Data().TotalCount)
	assert.Equal(t, int64(1), st.Data().ResponseCounts["/a"])
	assert.Equal(t, "ok\n", string(ctx.Response.Body()))

}
