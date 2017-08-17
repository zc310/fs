package main

import (
	"flag"
	"fmt"
	"log"

	"time"

	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":7070", "TCP address to listen to")
)

func main() {
	flag.Parse()
	if err := fasthttp.ListenAndServe(*addr, requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n\n")
	fmt.Fprintf(ctx, "%s", &ctx.Request)
	delay := string(ctx.QueryArgs().Peek("delay"))
	if delay != "" {
		d, err := time.ParseDuration(delay)
		if err != nil {
			fmt.Fprintln(ctx, err)
			return
		}
		time.Sleep(d)
	}
	fmt.Println(time.Now(), ctx.Request.URI().String())
}
