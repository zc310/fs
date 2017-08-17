package middleware

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/zc310/log"
)

type Plugin interface {
	Init(c *Config) error
	UnInit()
	Process(fasthttp.RequestHandler) fasthttp.RequestHandler
	Handler() fasthttp.RequestHandler
}

type Config struct {
	Logger log.Logger
	Router *fasthttprouter.Router
	Path   string
}
type NewFunc func() Plugin

var (
	SupportedPlugins = make(map[string]NewFunc)
)

func init() {
	SupportedPlugins["ok"] = func() Plugin { return &Ok{} }
	SupportedPlugins["favicon"] = func() Plugin { return &Favicon{} }
	SupportedPlugins["helloworld"] = func() Plugin { return &Helloworld{} }
	SupportedPlugins["compress"] = func() Plugin { return &Compress{} }
	SupportedPlugins["recover"] = func() Plugin { return &Recover{} }
	SupportedPlugins["log"] = func() Plugin { return &Logger{} }
	SupportedPlugins["pprof"] = func() Plugin { return &Pprof{} }
	SupportedPlugins["expvar"] = func() Plugin { return &ExpVar{} }
	SupportedPlugins["log"] = func() Plugin { return &Logger{} }
	SupportedPlugins["static"] = func() Plugin { return &Static{} }
	SupportedPlugins["stats"] = func() Plugin { return &Stats{} }
	SupportedPlugins["header"] = func() Plugin { return &Header{} }
	SupportedPlugins["proxy_header"] = func() Plugin { return &ProxyHeader{} }
	SupportedPlugins["proxy"] = func() Plugin { return &Proxy{} }
	SupportedPlugins["ratelimit"] = func() Plugin { return &Ratelimit{} }
	SupportedPlugins["cache"] = func() Plugin { return &Cache{} }
	SupportedPlugins["delay"] = func() Plugin { return &Delay{} }

	SupportedPlugins["fastcgi"] = func() Plugin { return &Fastcgi{} }
	SupportedPlugins["limit"] = func() Plugin { return &Limit{} }
	SupportedPlugins["basicauth"] = func() Plugin { return &BasicAuth{} }
	SupportedPlugins["expires"] = func() Plugin { return &Expires{} }
	SupportedPlugins["gallery"] = func() Plugin { return &Gallery{} }
}

func RegisterPlugin(name string, plugin NewFunc) {
	if name == "" {
		panic("plugin must have a name")
	}
	if _, ok := SupportedPlugins[name]; ok {
		panic("plugin named " + name + " already registered")
	}
	SupportedPlugins[name] = plugin
}
