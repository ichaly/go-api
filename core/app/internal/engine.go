package internal

import (
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewEngine(c *serv.Config, s *chi.Mux) (svc *serv.Service, err error) {
	svc, err = serv.NewGraphJinService(c)
	return
}
