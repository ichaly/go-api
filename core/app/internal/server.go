package internal

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	return r
}
