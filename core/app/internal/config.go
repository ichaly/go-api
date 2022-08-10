package internal

import (
	"github.com/dosco/graphjin/serv"
	"github.com/ichaly/go-api/core/app/internal/util"
	"path"
)

type Engine = serv.Config

type Config struct {
	// Engine holds config values for the GraphJin compiler
	Engine `mapstructure:",squash"`
}

func NewConfig() (*Config, error) {
	c, err := serv.ReadInConfig(path.Join(util.Root(), "./config", serv.GetConfigName()))
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = c.GetViper().Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	cfg.Engine = *c
	return cfg, nil
}
