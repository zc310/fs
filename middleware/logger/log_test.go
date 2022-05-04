package logger_test

import (
	"fmt"

	"github.com/valyala/fasthttp"

	. "github.com/zc310/fs/middleware/logger"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestLogHandler(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.Header.SetHost("a.com")
	req.Header.SetMethod("POST")
	req.Header.Set("User-Agent", "IE10")
	req.Header.Set("Referer", "http://b.com/")
	req.SetRequestURI("/a/b/c?asdfasdfas=asdfasdsa+d+f%09%27%27dd")
	ctx.Init(&req, nil, nil)

	fmt.Fprint(&ctx, "123456789")
	BuildCommonLogLine(os.Stdout, &ctx, time.Now())
}

func BenchmarkLogHandler(b *testing.B) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.Header.SetHost("www.google.com")
	req.Header.SetMethod("POST")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36")
	req.Header.Set("Referer", "http://www.baidu.com/")
	req.SetRequestURI("/a/b/c")
	ctx.Init(&req, nil, nil)

	fmt.Fprint(&ctx, "123456789")

	for n := 0; n < b.N; n++ {
		BuildCommonLogLine(ioutil.Discard, &ctx, time.Now())
	}
}
