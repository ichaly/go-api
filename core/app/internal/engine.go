package internal

import (
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewEngine(c *serv.Config, r *chi.Mux) (svc *serv.Service, err error) {
	if svc, err = serv.NewGraphJinService(c); err != nil {
		return
	}
	err = svc.Attach(r)
	return
}
