package middleware

import (
	"encoding/json"
	"fmt"

	"os"

	"github.com/valyala/fasthttp"

	"strings"
	"sync"
	"time"
)

type Stats struct {
	mu             sync.RWMutex     `json:"-"`
	Uptime         time.Time        `json:"-"`
	Pid            int              `json:"-"`
	ResponseCounts map[string]int64 `json:"-"`
	Path           string           `json:"path"`
}

type Data struct {
	UpTime         string           `json:"uptime"`
	UpTimeSec      float64          `json:"uptime_sec"`
	ResponseCounts map[string]int64 `json:"count"`
	TotalCount     int64            `json:"total_count"`
}

func (p *Stats) Init(c *Config) error {
	p.Uptime = time.Now()
	p.ResponseCounts = map[string]int64{}
	p.Pid = os.Getegid()

	if p.Path == "" {
		p.Path = "/_stats"
	}
	if !strings.HasPrefix(p.Path, "/") {
		p.Path = "/" + p.Path
	}

	c.Router.GET(p.Path, p.Handler())
	return nil
}
func (p *Stats) UnInit() {}
func (p *Stats) Data() *Data {
	now := time.Now()
	p.mu.RLock()
	uptime := now.Sub(p.Uptime)
	r := &Data{}
	count := int64(0)
	responseCounts := make(map[string]int64)
	for path, current := range p.ResponseCounts {
		responseCounts[path] = current
		count += current
	}
	p.mu.RUnlock()

	r.UpTime = uptime.String()
	r.UpTimeSec = uptime.Seconds()
	r.TotalCount = count
	r.ResponseCounts = responseCounts
	return r
}

func (p *Stats) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
		p.mu.Lock()
		p.ResponseCounts[string(ctx.Path())]++
		p.mu.Unlock()
	}
}

func (p *Stats) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json; charset=utf-8")
		err := json.NewEncoder(ctx.Response.BodyWriter()).Encode(p.Data())
		if err != nil {
			fmt.Fprint(ctx, "{}\n")
		}
	}

}
