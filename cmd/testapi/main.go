package main

import (
	"path/filepath"

	"github.com/valyala/fasthttp"
	"github.com/zc310/fasthttprouter"
	"github.com/zc310/fs/middleware"
	"github.com/zc310/log"
	"github.com/zc310/utils/fasthttputil"
	"github.com/zc310/utils/fasthttputil/favicon"

	"github.com/zc310/alice"

	"encoding/json"
	"flag"
)

type CompareResult struct {
	Code int         `json:"code,string"`
	Data interface{} `json:"data"`
}

var (
	port      = flag.String("port", ":7000", "port")
	client360 *fasthttp.HostClient
)

func main() {
	flag.Parse()
	log.SetPath("logs/")
	router := fasthttprouter.New()
	dir := ""
	var err error
	client360 = &fasthttp.HostClient{Addr: "cp.360.cn", MaxConns: 50}

	cfg := &middleware.Config{log.NewWithPrefix("proxy"), router, "/"}
	mwlog := &middleware.Logger{Filename: dir + "/logs/access.log",
		MaxSize:    30,
		MaxBackups: 7,
		MaxAge:     7,
		Compress:   false}
	mwlog.Init(cfg)
	gz := &middleware.Compress{}
	dir = filepath.Join(dir, "cache")

	mwstat := &middleware.Stats{}
	err = mwstat.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	mwcache := &middleware.Cache{}
	mwcache.Store = map[string]interface{}{"name": "file", "path": dir}
	mwcache.Key = "{document_root}.{request_body}"
	mwcache.Hash = "sha1"
	mwcache.Timeout = "12h"
	err = mwcache.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}
	sf := middleware.Singleflight{}
	sf.Init(cfg)
	sf.Key = "{document_root}.{request_body}"
	sf.Hash = "sha1"

	router.GET("/", fasthttputil.Ok)
	router.GET("/favicon.ico", favicon.Handler)

	mw := []alice.Constructor{mwlog.Process, mwstat.Process, mwcache.Process, gz.Process, sf.Process}
	//mw := []alice.Constructor{mwlog.Process, mwstat.Process,  gz.Process, sf.Process}

	router.POST("/tools/getCompareResult", alice.New(mw...).Then(api))

	server := &fasthttp.Server{
		MaxRequestBodySize: 70 * 1024 * 1024,
		Name:               "nginx/1.0",
		Handler:            router.Handler,
		Logger:             fasthttp.Logger(log.NewWithPrefix("default")),
	}

	log.Fatal(server.ListenAndServe(*port))
}

func api(ctx *fasthttp.RequestCtx) {
	var err error
	cr := &CompareResult{}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("http://cp.360.cn/tools/getCompareResult")
	req.Header.SetMethodBytes([]byte("POST"))
	req.Header.Set("Host", "cp.360.cn")
	req.Header.SetUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept-Encoding", string(ctx.Request.Header.Peek("Accept-Encoding")))
	req.Header.Del("Connection")

	req.SetBody(ctx.Request.Body())

	//fmt.Println("###",req.Header.String())
	err = client360.Do(req, &ctx.Response)
	//fmt.Println("###",ctx.Response.Header.String())
	ctx.Response.Header.Del("Connection")

	ctx.Response.Header.Del("Cache-Control")
	ctx.Response.Header.Del("Expires")
	ctx.Response.Header.Del("max-age")
	ctx.Response.Header.Del("Pragma")

	err = json.Unmarshal(ctx.Response.Body(), cr)
	if (err != nil) || (cr.Code != 200) {
		ctx.Response.Header.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		ctx.Response.Header.Set("Pragma", "no-cache")
		ctx.Response.Header.Set("Expires", "Fri, 29 Aug 1997 02:14:00 EST")

	}

}
