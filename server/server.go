package server

import (
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/valyala/fasthttp"
	"github.com/zc310/alice"
	"github.com/zc310/fasthttprouter"
	"github.com/zc310/fs/middleware"
	"github.com/zc310/fs/template"
	"github.com/zc310/log"
	"github.com/zc310/utils/fasthttputil"
)

// Start 初始化
func Start(c *Config) error {
	cfg := new(middleware.Config)
	cfg.Router = fasthttprouter.New()
	if c.Log.Path != "" {
		t, err := template.NewTemplate(c.Log.Path)
		if err != nil {
			return err
		}
		log.SetPath(string(t))
	}

	cfg.Logger = log.NewWithPrefix("fs")
	cfg.Path = "/"
	cfg.Logger.Print(c)
	var err error
	var mw []alice.Constructor
	var mp middleware.Plugin
	var h alice.Chain
	if len(c.Middleware) > 0 {
		mw, err = c.Middleware.Load(cfg)
		if err != nil {
			return err
		}
		h = alice.New(mw...)
	}
	hw := make(fasthttputil.HostSwitch)
	for _, h := range c.Handler {
		cfg.Router = fasthttprouter.New()
		mw, err = h.Middleware.Load(cfg)
		if err != nil {
			return err
		}
		for _, name := range h.Host {
			hw.Add(name, alice.New(mw...).Then(cfg.Router.Handler))
		}

		for _, router := range h.Router {
			mw, err = router.Middleware.Load(cfg)
			if err != nil {
				return err
			}
			mp, err = router.Handler.Load(cfg)
			if err != nil {
				return err
			}
			for _, path := range router.Paths {
				cfg.Path = strings.TrimSuffix(path, "*filepath")
				AddRouter(cfg.Router, path, alice.New(mw...).Then(mp.Handler()))
			}

		}
	}
	m, err := humanize.ParseBytes(c.MaxBodySize)
	if err != nil {
		cfg.Logger.Print(err)
		m = fasthttp.DefaultMaxRequestBodySize
	}
	cfg.Logger.Print("DefaultMaxRequestBodySize", c.MaxBodySize, m)
	server := &fasthttp.Server{
		Name:               c.Name,
		Handler:            h.Then(hw.Handler),
		Logger:             cfg.Logger,
		MaxRequestBodySize: int(m),
	}
	return server.ListenAndServe(c.Listen)
}
func AddRouter(router *fasthttprouter.Router, path string, handle fasthttp.RequestHandler) {
	router.GET(path, handle)
	router.POST(path, handle)
	router.HEAD(path, handle)
	router.OPTIONS(path, handle)
	router.DELETE(path, handle)
	router.PUT(path, handle)
	router.PATCH(path, handle)
}
