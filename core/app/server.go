package app

import (
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewServer(c *serv.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	return r
}
