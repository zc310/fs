package server

import (
	"encoding/json"
	"io/ioutil"

	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/zc310/alice"
	"github.com/zc310/apiproxy/middleware"

	"github.com/zc310/log"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"strings"
)

type (
	Plugin     map[string]interface{}
	Middleware []Plugin
	Host       struct {
		Middleware Middleware       `json:"middleware"`
		Paths      map[string]*Path `json:"paths"`
	}

	Path struct {
		Middleware Middleware `json:"middleware"`
		Handler    Plugin     `json:"handler"`
	}

	// Config  配置
	Config struct {
		Name string `json:"name"`
		// listen 绑定端口
		Listen      string     `json:"listen"`
		MaxBodySize string     `json:"max_body_size"`
		Middleware  Middleware `json:"middleware"`
		// Hosts Host列表
		Hosts map[string]*Host `json:"hosts"`
		Log   struct {
			Path string `json:"path"`
		} `json:"log"`
		Logger log.Logger `json:"-",yaml:"-"`
	}
)

// ReadFile 读取配置文件
func (p *Config) ReadFile(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if strings.ToLower(filepath.Ext(filename)) == ".yaml" {
		return yaml.Unmarshal(b, p)
	}
	return json.Unmarshal(b, p)
}
func (p *Plugin) Load(c *middleware.Config) (middleware.Plugin, error) {
	name := (*p)["name"].(string)
	pf, ok := middleware.SupportedPlugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	pi := pf()
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  pi,
	})
	if err != nil {
		return nil, err
	}

	if err = dec.Decode(p); err != nil {
		return nil, err
	}
	err = pi.Init(c)
	if err != nil {
		return pi, err
	}
	return pi, nil
}
func (p Middleware) Load(c *middleware.Config) ([]alice.Constructor, error) {
	var pl []alice.Constructor
	var pi middleware.Plugin
	var err error
	for _, ps := range p {
		pi, err = ps.Load(c)
		if err != nil {
			return pl, err
		}
		pl = append(pl, pi.Process)
	}
	return pl, err
}
