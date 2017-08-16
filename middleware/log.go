package middleware

import (
	"io"
	"net"
	"strconv"
	"time"

	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"github.com/zc310/apiproxy/template"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Filename   string `json:"filename"`
	MaxSize    int    `json:"rotate_size"`
	MaxAge     int    `json:"rotate_age"`
	MaxBackups int    `json:"rotate_backups"`
	Compress   bool   `json:"rotate_compress"`
	Log        io.Writer
}

func (p *Logger) Init(c *Config) (err error) {
	filename, err := template.NewTemplate(p.Filename)
	if err != nil {
		return err
	}
	p.Log = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    p.MaxSize,
		MaxBackups: p.MaxBackups,
		MaxAge:     p.MaxAge,
		Compress:   p.Compress,
	}

	return nil
}
func (p *Logger) UnInit() {}
func (p *Logger) Process(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		t := time.Now()
		h(ctx)
		buildCommonLogLine(p.Log, ctx, t)
	}
}
func (p *Logger) Handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {}
}
func buildCommonLogLine(w io.Writer, ctx *fasthttp.RequestCtx, ts time.Time) {
	host, _, err := net.SplitHostPort(ctx.RemoteAddr().String())
	if err != nil {
		host = ctx.RemoteAddr().String()
	}
	buf := bytebufferpool.Get()
	buf.WriteString(host)
	buf.WriteString(" - ")
	buf.WriteString("-")
	buf.WriteString(" [")
	buf.WriteString(ts.Format("02/Jan/2006:15:04:05 -0700"))
	buf.WriteString(`] "`)
	buf.Write(ctx.Method())
	buf.WriteString(" ")
	buf.WriteString(strconv.Quote(string(ctx.RequestURI())))
	buf.WriteString(" ")
	buf.WriteString("HTTP/1.0")
	buf.WriteString(`" `)
	buf.WriteString(strconv.Itoa(ctx.Response.StatusCode()))
	buf.WriteString(" ")
	buf.WriteString(strconv.Itoa(len(ctx.Response.Body())))
	buf.WriteString(` "`)
	buf.Write(ctx.Referer())
	buf.WriteString(`" "`)
	buf.Write(ctx.UserAgent())
	buf.WriteString(`"`)
	buf.WriteString("\n")
	w.Write(buf.B)
	bytebufferpool.Put(buf)

}
