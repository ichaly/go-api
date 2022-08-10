package internal

import (
	"fmt"
	"github.com/dosco/graphjin/core"
	"github.com/ichaly/go-api/core/app/internal/util"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func NewViper() (*viper.Viper, error) {
	root := path.Join(util.Root(), "./config")
	vi := newViper(root, core.GetConfigName())
	err := vi.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if pcf := vi.GetString("inherits"); pcf != "" {
		cf := vi.ConfigFileUsed()
		vi = newViper(root, pcf)

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

	return vi, nil
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
	vi.SetDefault("database.user", "postgres")
	vi.SetDefault("database.password", "")
	vi.SetDefault("database.schema", "public")
	vi.SetDefault("database.pool_size", 10)

	vi.SetDefault("env", "development")

	_ = vi.BindEnv("env", "GO_ENV") //nolint: errcheck
	_ = vi.BindEnv("host", "HOST")  //nolint: errcheck
	_ = vi.BindEnv("port", "PORT")  //nolint: errcheck

	vi.SetDefault("auth.rails.max_idle", 80)
	vi.SetDefault("auth.rails.max_active", 12000)
	vi.SetDefault("auth.subs_creds_in_vars", false)

	return vi
}
