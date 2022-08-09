package internal

import (
	"github.com/dosco/graphjin/serv"
	"github.com/ichaly/go-api/core/app/internal/util"
	"path"
)

func NewConfig() (*serv.Config, error) {
	return serv.ReadInConfig(path.Join(util.Root(), "./config", serv.GetConfigName()))
}
