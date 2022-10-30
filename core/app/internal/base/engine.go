package base

import (
	"github.com/dosco/graphjin/serv"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewEngine(c *Config) (*serv.Service, error) {
	return serv.NewGraphJinService(&c.Engine)
}
