package internal

import (
	"github.com/dosco/graphjin/core"
	"github.com/spf13/viper"
)

type Core = core.Config

type Database struct {
	Type     string     `mapstructure:"type"`
	Url      string     `mapstructure:"url"`
	Host     string     `mapstructure:"host"`
	Port     int        `mapstructure:"port"`
	Name     string     `mapstructure:"dbname"`
	Username string     `mapstructure:"user"`
	Password string     `mapstructure:"password"`
	Schema   string     `mapstructure:"schema"`
	Sources  []Database `mapstructure:"sources"`
	Replicas []Database `mapstructure:"replicas"`
}

type Config struct {
	// Core holds config values for the GraphJin compiler
	Core `mapstructure:",squash"`
	Database
}

func NewConfig(vi *viper.Viper) (*Config, error) {
	cfg := &Config{}
	err := vi.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
