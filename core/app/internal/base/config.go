package base

import (
	"github.com/dosco/graphjin/serv"
	"github.com/ichaly/go-api/core/app/internal/util"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"path"
)

type Engine = serv.Config

type Config struct {
	// Engine holds config values for the GraphJin compiler
	Engine `mapstructure:",squash"`
	Driver *base64Captcha.DriverString `mapstructure:"captcha"`
}

func NewConfig() (*Config, error) {
	c, err := serv.ReadInConfig(path.Join(util.Root(), "./config", serv.GetConfigName()))
	if err != nil {
		return nil, err
	}

	v := c.GetViper()
	v.RegisterAlias("captcha.BgColor", "captcha.bg-color")
	v.RegisterAlias("captcha.NoiseCount", "captcha.noise-count")

	v.SetDefault("captcha.BgColor", color.RGBA{A: 255, R: 233, G: 238, B: 243})
	v.SetDefault("captcha.NoiseCount", 20)
	v.SetDefault("captcha.Fonts", []string{"3Dumb.ttf"})

	cfg := &Config{}
	err = v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	cfg.Engine = *c
	return cfg, nil
}
