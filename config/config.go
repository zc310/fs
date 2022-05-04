package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/zc310/fs/middleware"
	"github.com/zc310/fs/template"

	"github.com/mitchellh/mapstructure"

	"github.com/zc310/alice"

	"github.com/zc310/log"
	"gopkg.in/yaml.v3"
)

type (
	Plugin     map[string]interface{}
	Middleware []Plugin
	Handler    struct {
		Host       []string   `json:"host"`
		Middleware Middleware `json:"middleware"`
		Router     []*Router  `json:"router"`
	}

	Router struct {
		Paths      []string   `json:"paths"`
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
		// Handler Handler
		Handler []*Handler `json:"handler"`
		Log     struct {
			Path string `json:"path"`
		} `json:"log"`
		Logger log.Logger `json:"-" yaml:"-"`
	}
	// CheckList key list
	CheckList map[string][]string
)

// CheckHit check args
func CheckHit(list CheckList, tpl *template.Template) (bool, error) {
	if list != nil {
		var key []byte
		var t string
		var err error
		for k, v1 := range list {
			if key, err = tpl.Execute(k); err != nil {
				return false, err
			}
			t = string(key)
			for _, v2 := range v1 {
				if v2 == t {
					return true, nil
				}
			}
		}
		return false, nil
	}
	return true, nil
}

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
