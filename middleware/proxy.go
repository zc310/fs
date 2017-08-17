package middleware

import (
	"github.com/valyala/fasthttp"

	"github.com/zc310/log"
	"github.com/zc310/utils/fasthttputil"

	"errors"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/zc310/utils"
)

type Proxy struct {
	Balance map[string]interface{} `json:"balance"`
	Targets []*fasthttputil.Target `json:"targets"`
	pool    HostPool
	policy  Policy
	log     log.Logger
	stop    chan struct{}
}
type Balance struct {
	Name string `json:"name"`
}

func (p *Proxy) Init(c *Config) (err error) {
	p.log = c.Logger.NewWithPrefix("proxy")
	var py *fasthttputil.Proxy
	for _, to := range p.Targets {
		py, err = fasthttputil.NewProxyClient(to)
		if err != nil {
			return
		}
		p.pool = append(p.pool, py)

	}
	err = p.loadBalance()
	if err != nil {
		return err
	}
	p.stop = make(chan struct{})
	go p.Cron()
	return
}
func (p *Proxy) UnInit() { close(p.stop) }
func (p *Proxy) loadBalance() error {
	name := utils.GetString(p.Balance["name"])

	pf, ok := supportedPolicies[name]
	if ok {
		p.policy = pf()
	} else {
		p.policy = supportedPolicies["random"]()
	}
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  p.policy,
	})
	if err != nil {
		return err
	}

	if err = dec.Decode(p.Balance); err != nil {
		return err
	}

	return p.policy.Init()
}
func (p *Proxy) do(ctx *fasthttp.RequestCtx) error {
	proxy := p.policy.Select(p.pool, ctx)
	if proxy == nil {
		return errors.New("proxy for host '" + string(ctx.Host()) + "' is nil")

	}
	return proxy.Handler(ctx)
}
func (p *Proxy) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		const maxAttempts = 3
		attempts := 0
		for {
			err := p.do(ctx)
			if err == nil {
				break
			}
			p.log.Error(err)
			attempts++
			if attempts > maxAttempts {
				break
			}
		}
	}
}
func (p *Proxy) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}
}
func (p *Proxy) healthCheck() {
	for _, host := range p.pool {
		host.HealthCheck()
	}
}
func (p *Proxy) Cron() {
	ticker := time.NewTicker(time.Second * 10)
	p.healthCheck()
	for {
		select {
		case <-ticker.C:
			p.healthCheck()
		case <-p.stop:
			ticker.Stop()
			return
		}
	}
}
