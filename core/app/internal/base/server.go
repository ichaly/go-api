package base

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"time"
)

func NewServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	return r
}
