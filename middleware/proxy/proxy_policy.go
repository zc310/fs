package proxy

import (
	"hash/fnv"

	"github.com/valyala/fasthttp"
	"github.com/zc310/utils/fasthttputil"

	"net"

	"github.com/valyala/bytebufferpool"
	"math/rand"
	"sync"
)

var (
	supportedPolicies = make(map[string]func() Policy)
)

type HostPool []*fasthttputil.Proxy
type Policy interface {
	Select(pool HostPool, ctx *fasthttp.RequestCtx) *fasthttputil.Proxy
	Init() error
}

func init() {
	RegisterPolicy("hash_ip", func() Policy { return &IPHash{} })
	RegisterPolicy("hash_body", func() Policy { return &BodyHash{} })
	RegisterPolicy("first", func() Policy { return &First{} })
	RegisterPolicy("random", func() Policy { return &Random{} })
	RegisterPolicy("round_robin", func() Policy { return &RoundRobin{} })
}
func RegisterPolicy(name string, policy func() Policy) {
	supportedPolicies[name] = policy
}
func Fnv(b []byte) uint32 {
	h := fnv.New32a()
	h.Write(b)
	return h.Sum32()
}

// IPHash is a policy that selects hosts based on hashing the request ip
type IPHash struct{}

func (r *IPHash) Init() error { return nil }

func (r *IPHash) Select(pool HostPool, ctx *fasthttp.RequestCtx) *fasthttputil.Proxy {
	n := uint32(len(pool))
	clientIP, _, err := net.SplitHostPort(ctx.RemoteAddr().String())
	if err != nil {
		clientIP = ctx.RemoteAddr().String()
	}
	index := Fnv([]byte(clientIP)) % n
	for i := uint32(0); i < n; i++ {
		index += i
		host := pool[index%n]
		if host.Available() {
			return host
		}
	}
	return nil
}

// BodyHash is a policy that selects hosts based on hashing the request args
type BodyHash struct {
	Args       []string `json:"args"`
	bufferpool bytebufferpool.Pool
}

func (r *BodyHash) Init() error { return nil }
func (r *BodyHash) Select(pool HostPool, ctx *fasthttp.RequestCtx) *fasthttputil.Proxy {
	n := uint32(len(pool))
	var index uint32

	if len(r.Args) > 0 {
		buf := r.bufferpool.Get()
		for _, v := range r.Args {
			buf.Write(fasthttputil.GetArgs(ctx, v))
			buf.Write([]byte("\n"))
		}

		index = Fnv(buf.B) % n
		r.bufferpool.Put(buf)
	} else {
		index = Fnv(ctx.Request.Body()) % n
	}

	for i := uint32(0); i < n; i++ {
		index += i
		host := pool[index%n]
		if host.Available() {
			return host
		}
	}

	return nil
}

// First is a policy that selects hosts based front active host
type First struct{}

func (r *First) Init() error { return nil }
func (r *First) Select(pool HostPool, ctx *fasthttp.RequestCtx) *fasthttputil.Proxy {
	for i := 0; i < len(pool); i++ {
		host := pool[i]
		if host.Available() {
			return host
		}
	}
	return nil
}

// Random is a policy that selects up hosts from a pool at random.
type Random struct{}

func (r *Random) Init() error { return nil }

// Select selects an up host at random from the specified pool.
func (r *Random) Select(pool HostPool, ctx *fasthttp.RequestCtx) *fasthttputil.Proxy {
	// Because the number of available hosts isn't known
	// up front, the host is selected via reservoir sampling
	// https://en.wikipedia.org/wiki/Reservoir_sampling
	var randHost *fasthttputil.Proxy
	count := 0
	for _, host := range pool {
		if !host.Available() {
			continue
		}

		// (n % 1 == 0) holds for all n, therefore randHost
		// will always get assigned a value if there is
		// at least 1 available host
		count++
		if (rand.Int() % count) == 0 {
			randHost = host
		}

	}
	return randHost
}

// RoundRobin is a policy that selects hosts based on round robin ordering.
type RoundRobin struct {
	robin uint32
	mutex sync.Mutex
}

func (r *RoundRobin) Init() error { return nil }

// Select selects an up host from the pool using a round robin ordering scheme.
func (r *RoundRobin) Select(pool HostPool, ctx *fasthttp.RequestCtx) *fasthttputil.Proxy {
	poolLen := uint32(len(pool))
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// Return next available host
	for i := uint32(0); i < poolLen; i++ {
		r.robin++
		host := pool[r.robin%poolLen]
		if host.Available() {
			return host
		}
	}
	return nil
}
