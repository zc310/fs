package middleware

import (
	"github.com/valyala/fasthttp"
	"github.com/zc310/fasthttprouter"
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

func RegisterPlugin(name string, plugin NewFunc) {
	if name == "" {
		panic("plugin must have a name")
	}
	if _, ok := SupportedPlugins[name]; ok {
		panic("plugin named " + name + " already registered")
	}
	SupportedPlugins[name] = plugin
}
