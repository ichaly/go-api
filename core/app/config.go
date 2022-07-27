package app

import (
	"github.com/dosco/graphjin/core"
	"github.com/spf13/viper"
	"path"
	"path/filepath"
)

type (
	Core = core.Config

	Database struct {
		Type     string
		Url      string
		Host     string
		Port     int
		Name     string
		Username string
		Password string
		Sources  []Database
		Replicas []Database
	}

	Config struct {
		Endpoint string

		// Core holds config values for the GraphJin compiler
		Core `mapstructure:",squash"`

		DataSource Database
		CacheStore Database

		// SecretsFile is the file for the secret key store
		SecretsFile string `mapstructure:"secrets_file"`
		// RateLimiter sets the API rate limits
		RateLimiter RateLimiter `mapstructure:"rate_limiter"`

		secrets map[string]string
		hash    string
		name    string
		vi      *viper.Viper
	}

	// RateLimiter sets the API rate limits
	RateLimiter struct {
		Rate     float64
		Bucket   int
		IPHeader string `mapstructure:"ip_header"`
	}
)

func NewConfig() (cfg *Config, err error) {
	root, err := filepath.Abs("./config")
	if err != nil {
		return
	}
	cfg, err = ReadInConfig(path.Join(root, core.GetConfigName()))
	if err != nil {
		return
	}
	return
}

func ReadInConfig(configFile string) (cfg *Config, err error) {
	return
}
