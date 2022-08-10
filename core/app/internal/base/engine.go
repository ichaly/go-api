package base

import (
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewEngine(c *Config, r *chi.Mux) (*serv.Service, error) {
	svc, err := serv.NewGraphJinService(&c.Engine)
	if err != nil {
		return nil, err
	}
	err = svc.Attach(r)
	if err != nil {
		return nil, err
	}
	return svc, nil
}
