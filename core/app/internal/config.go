package internal

import (
	"github.com/dosco/graphjin/serv"
	"path"
	"path/filepath"
)

func NewConfig() (cfg *serv.Config, err error) {
	root, err := filepath.Abs("./config")
	if err != nil {
		return
	}
	cfg, err = serv.ReadInConfig(path.Join(root, serv.GetConfigName()))
	return
}
