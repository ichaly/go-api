package base

import (
	"fmt"
	"github.com/dosco/graphjin/core"
	"github.com/ichaly/go-api/core/app/pkg/util"
	"github.com/json-iterator/go"
	"github.com/spf13/viper"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Core = core.Config

type Database struct {
	Type     string `jsonschema:"title=Type,enum=postgres,enum=mysql,enum=redis,enum=memory"`
	Host     string `jsonschema:"title=Host"`
	Port     uint16 `jsonschema:"title=Port"`
	Name     string `jsonschema:"title=Database Name"`
	Username string `jsonschema:"title=Username"`
	Password string `jsonschema:"title=Password"`
}

type Config struct {
	vi *viper.Viper

	Core     `mapstructure:",squash"`
	Cache    *Database `mapstructure:"cache" jsonschema:"title=Cache"`
	Database *Database `mapstructure:"database" jsonschema:"title=Database"`

	HostPort string `mapstructure:"host_port" jsonschema:"title=Host and Port"`
	AppName  string `mapstructure:"app_name" jsonschema:"title=Application Name"`
}

func NewConfig() (*Config, error) {
	return readInConfig(path.Join(util.Root(), "./config", core.GetConfigName()))
}

func readInConfig(configFile string) (*Config, error) {
	cp := filepath.Dir(configFile)
	vi := newViper(cp, filepath.Base(configFile))

	if err := vi.ReadInConfig(); err != nil {
		return nil, err
	}

	if pcf := vi.GetString("inherits"); pcf != "" {
		cf := vi.ConfigFileUsed()
		vi = newViper(cp, pcf)

		if err := vi.ReadInConfig(); err != nil {
			return nil, err
		}

		if v := vi.GetString("inherits"); v != "" {
			return nil, fmt.Errorf("inherited config (%s) cannot itself inherit (%s)", pcf, v)
		}

		vi.SetConfigFile(cf)

		if err := vi.MergeInConfig(); err != nil {
			return nil, err
		}
	}

	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GJ_") || strings.HasPrefix(e, "SJ_") {
			kv := strings.SplitN(e, "=", 2)
			util.SetKeyValue(vi, kv[0], kv[1])
		}
	}

	c := &Config{vi: vi}
	c.ConfigPath = cp

	if err := vi.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to decode config, %v", err)
	}

	return c, nil
}

func newViper(configPath, configFile string) *viper.Viper {
	vi := newViperWithDefaults()
	vi.SetConfigName(strings.TrimSuffix(configFile, filepath.Ext(configFile)))

	if configPath == "" {
		vi.AddConfigPath("./config")
	} else {
		vi.AddConfigPath(configPath)
	}

	return vi
}

func newViperWithDefaults() *viper.Viper {
	vi := viper.New()

	vi.SetDefault("host_port", "0.0.0.0:8080")
	vi.SetDefault("web_ui", false)
	vi.SetDefault("enable_tracing", false)
	vi.SetDefault("auth_fail_block", false)
	vi.SetDefault("seed_file", "seed.js")

	vi.SetDefault("log_level", "info")
	vi.SetDefault("log_format", "json")

	vi.SetDefault("default_block", true)

	vi.SetDefault("database.type", "postgres")
	vi.SetDefault("database.host", "localhost")
	vi.SetDefault("database.port", 5432)
	vi.SetDefault("database.username", "postgres")
	vi.SetDefault("database.password", "")
	vi.SetDefault("database.schema", "public")
	vi.SetDefault("database.pool_size", 10)

	vi.SetDefault("env", "development")

	_ = vi.BindEnv("env", "GO_ENV")
	_ = vi.BindEnv("host", "HOST")
	_ = vi.BindEnv("port", "PORT")

	vi.SetDefault("auth.rails.max_idle", 80)
	vi.SetDefault("auth.rails.max_active", 12000)
	vi.SetDefault("auth.subs_creds_in_vars", false)

	vi.RegisterAlias("captcha.BgColor", "captcha.bg-color")
	vi.RegisterAlias("captcha.NoiseCount", "captcha.noise-count")

	vi.SetDefault("cache.type", "memory")
	vi.SetDefault("captcha.BgColor", color.RGBA{A: 255, R: 233, G: 238, B: 243})
	vi.SetDefault("captcha.NoiseCount", 20)
	vi.SetDefault("captcha.Fonts", []string{"3Dumb.ttf"})

	return vi
}
