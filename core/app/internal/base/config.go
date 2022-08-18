package base

import (
	"github.com/dosco/graphjin/serv"
	"github.com/ichaly/go-api/core/app/internal/util"
	"github.com/mojocn/base64Captcha"
	"path"
)

type Engine = serv.Config

type Config struct {
	// Engine holds config values for the GraphJin compiler
	Engine `mapstructure:",squash"`
	Driver *base64Captcha.DriverDigit `mapstructure:"captcha"`
}

func NewConfig() (*Config, error) {
	c, err := serv.ReadInConfig(path.Join(util.Root(), "./config", serv.GetConfigName()))
	if err != nil {
		return nil, err
	}

	v := c.GetViper()
	v.RegisterAlias("captcha.MaxSkew", "captcha.max-skew")
	v.RegisterAlias("captcha.DotCount", "captcha.dot-count")

	cfg := &Config{}
	err = v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	cfg.Engine = *c
	return cfg, nil
}
