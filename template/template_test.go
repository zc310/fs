package template

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestTemplate(t *testing.T) {
	s, err := New(nil).Execute("{tempdir}/logs/")
	assert.Equal(t, nil, err)
	assert.Equal(t, string(s), os.TempDir()+"/logs/")
}

func testExecute(t *testing.T, template *Template, s, out string) {
	output, err := template.Execute(s)
	assert.Equal(t, err, nil)
	assert.Equal(t, output, []byte(out))
}
func TestTemplate_Execute(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.Header.SetHost("a.com")
	req.Header.SetMethod("POST")
	req.Header.Set("User-Agent", "IE10")
	req.Header.Set("Content-Type", "image/svg+xml")
	req.Header.Set("Referer", "http://b.com")
	req.SetRequestURI("/a/b/c?asdfasdfas=asdfasdsa+d+f%09%27%27dd&a1=123")
	ctx.Init(&req, nil, nil)
	fmt.Fprint(ctx.Request.BodyWriter(), "123456789")

	tmp := New(nil)
	testExecute(t, tmp, "{tempdir}", os.TempDir())

	tmp = New(&ctx)
	testExecute(t, tmp, "{scheme}", "http")
	testExecute(t, tmp, "{host}", "a.com")
	testExecute(t, tmp, "{request_uri}", "http://a.com/a/b/c?asdfasdfas=asdfasdsa+d+f%09%27%27dd&a1=123")
	testExecute(t, tmp, "{request_method}", "POST")
	testExecute(t, tmp, "{request_body}", "123456789")
	testExecute(t, tmp, "{content_type}", "image/svg+xml")
	testExecute(t, tmp, "{content_length}", "9")
	testExecute(t, tmp, "{query_string}", "asdfasdfas=asdfasdsa+d+f%09%27%27dd&a1=123")
	testExecute(t, tmp, "{document_root}", "/a/b/c")
	testExecute(t, tmp, "{remote_addr}", "0.0.0.0")
	testExecute(t, tmp, "{a1}", "123")
	testExecute(t, tmp, "{User-Agent}", "IE10")
	testExecute(t, tmp, "{Referer}", "http://b.com")
}
