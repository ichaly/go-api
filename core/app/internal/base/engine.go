package base

import (
	"github.com/dosco/graphjin/serv"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewEngine(c *Config) (*serv.Service, error) {
	svc, err := serv.NewGraphJinService(&c.Engine)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return svc, nil
}
