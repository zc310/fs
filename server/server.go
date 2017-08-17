package server

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/dustin/go-humanize"
	"github.com/valyala/fasthttp"
	"github.com/zc310/alice"
	"github.com/zc310/fs/middleware"
	"github.com/zc310/fs/template"
	"github.com/zc310/log"
	"github.com/zc310/utils/fasthttputil"
	"strings"
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
		log.SetPath(t)
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
	for name, host := range c.Hosts {
		cfg.Router = fasthttprouter.New()
		mw, err = host.Middleware.Load(cfg)
		if err != nil {
			return err
		}

		hw.Add(name, alice.New(mw...).Then(cfg.Router.Handler))

		for path, p := range host.Paths {
			cfg.Path = strings.TrimSuffix(path, "*filepath")
			mw, err = p.Middleware.Load(cfg)
			if err != nil {
				return err
			}
			mp, err = p.Handler.Load(cfg)
			if err != nil {
				return err
			}
			AddRouter(cfg.Router, path, alice.New(mw...).Then(mp.Handler()))
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
