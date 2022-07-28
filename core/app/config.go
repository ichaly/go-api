package app

import (
	"database/sql"
	"github.com/dosco/graphjin/core"
	"github.com/spf13/viper"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type (
	Core = core.Config

	Database struct {
		Type            string
		Url             string
		Host            string
		Port            uint16
		Name            string `mapstructure:"  dbname"`
		Username        string `mapstructure:"user"`
		Password        string
		Schema          string
		PoolSize        int           `mapstructure:"pool_size"`
		MaxConnections  int           `mapstructure:"max_connections"`
		MaxConnIdleTime time.Duration `mapstructure:"max_connection_idle_time"`
		MaxConnLifeTime time.Duration `mapstructure:"max_connection_life_time"`
		PingTimeout     time.Duration `mapstructure:"ping_timeout"`
		EnableTLS       bool          `mapstructure:"enable_tls"`
		ServerName      string        `mapstructure:"server_name"`
		ServerCert      string        `mapstructure:"server_cert"`
		ClientCert      string        `mapstructure:"client_cert"`
		ClientKey       string        `mapstructure:"client_key"`
	}

	Config struct {
		// Core holds config values for the GraphJin compiler
		Core `mapstructure:",squash"`

		// AppName is the name of your application used in log and debug messages
		AppName string `mapstructure:"app_name"`

		// Production when set to true runs the service with production level
		// security and other defaults. For example allow lists are enforced.
		Production bool

		Endpoint string
		// HostPost to run the service on. Example localhost:8080
		HostPort string `mapstructure:"host_port"`
		// DataSource struct contains db config
		DataSource Database `mapstructure:"database"`
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
	// copy over db_type from database.type
	if cfg.Core.DBType == "" {
		cfg.Core.DBType = cfg.DataSource.Type
	}
	return
}

func ReadInConfig(configFile string) (cfg *Config, err error) {
	return
}

func (my *Config) Host() string {
	hp := strings.SplitN(my.HostPort, ":", 2)
	if len(hp) == 2 {
		return hp[0]
	}
	return "0.0.0.0"
}

func (my *Config) Port() string {
	hp := strings.SplitN(my.HostPort, ":", 2)
	if len(hp) == 2 {
		return hp[1]
	}
	return "8080"
}

func (my *Config) GetDB() (db *sql.DB) {
	return
}
