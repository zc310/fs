package plugin

import (
	. "github.com/zc310/fs/middleware"
	"github.com/zc310/fs/middleware/basicauth"
	"github.com/zc310/fs/middleware/cache"
	"github.com/zc310/fs/middleware/delay"
	"github.com/zc310/fs/middleware/expvar"
	"github.com/zc310/fs/middleware/favicon"
	"github.com/zc310/fs/middleware/gzip"
	"github.com/zc310/fs/middleware/header"
	"github.com/zc310/fs/middleware/logger"
	"github.com/zc310/fs/middleware/pprof"
	"github.com/zc310/fs/middleware/proxy"
	"github.com/zc310/fs/middleware/singleflight"
	"github.com/zc310/fs/middleware/static"
	"github.com/zc310/fs/middleware/stats"
)

func init() {
	SupportedPlugins["ok"] = func() Plugin { return &favicon.Ok{} }
	SupportedPlugins["favicon"] = func() Plugin { return &favicon.Favicon{} }
	SupportedPlugins["helloworld"] = func() Plugin { return &Helloworld{} }
	SupportedPlugins["compress"] = func() Plugin { return &gzip.Compress{} }
	SupportedPlugins["recover"] = func() Plugin { return &Recover{} }
	SupportedPlugins["log"] = func() Plugin { return &logger.Logger{} }
	SupportedPlugins["pprof"] = func() Plugin { return &pprof.Pprof{} }
	SupportedPlugins["expvar"] = func() Plugin { return &expvar.ExpVar{} }

	SupportedPlugins["static"] = func() Plugin { return &static.Static{} }
	SupportedPlugins["stats"] = func() Plugin { return &stats.Stats{} }
	SupportedPlugins["header"] = func() Plugin { return &header.Header{} }
	SupportedPlugins["proxy_header"] = func() Plugin { return &proxy.ProxyHeader{} }
	SupportedPlugins["proxy"] = func() Plugin { return &proxy.Proxy{} }
	SupportedPlugins["ratelimit"] = func() Plugin { return &proxy.Ratelimit{} }
	SupportedPlugins["cache"] = func() Plugin { return &cache.Cache{} }
	SupportedPlugins["delay"] = func() Plugin { return &delay.Delay{} }
	SupportedPlugins["singleflight"] = func() Plugin { return &singleflight.Singleflight{} }

	SupportedPlugins["fastcgi"] = func() Plugin { return &Fastcgi{} }
	SupportedPlugins["limit"] = func() Plugin { return &Limit{} }
	SupportedPlugins["basicauth"] = func() Plugin { return &basicauth.BasicAuth{} }
	SupportedPlugins["expires"] = func() Plugin { return &Expires{} }
	SupportedPlugins["gallery"] = func() Plugin { return &Gallery{} }
}
