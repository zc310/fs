package main

import (
	"flag"
	"github.com/zc310/apiproxy/server"
	"github.com/zc310/log"
)

type Proxy struct {
	Form   []string
	To     []string
	Name   string
	Cache  bool
	Policy string
}

func main() {
	var cfg string
	flag.StringVar(&cfg, "cfg", "fs.json", "fs.json or fs.yaml")
	flag.Parse()
	config := new(server.Config)
	err := config.ReadFile(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config)
	if err != nil {

		log.Fatal(err)
	}
}
