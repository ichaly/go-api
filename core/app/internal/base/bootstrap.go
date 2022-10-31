package base

import (
	"context"
	"fmt"
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi"
	"github.com/ichaly/go-api/core/app/pkg"
	"go.uber.org/fx"
	"net/http"
)

type Enhance struct {
	fx.In
	Plugins     []core.Plugin `group:"plugin"`
	Middlewares []core.Plugin `group:"middleware"`
}

func Bootstrap(
	l fx.Lifecycle, s *serv.Service, r *chi.Mux, c *Config, e Enhance,
) {
	//init middlewares
	for _, m := range e.Middlewares {
		if !m.Protected() {
			m.Init(r)
		}
	}
	//wrap graphql service
	r.Group(func(r chi.Router) {
		for _, m := range e.Middlewares {
			if m.Protected() {
				m.Init(r)
			}
		}
		_ = s.Attach(r)
		for _, p := range e.Plugins {
			if p.Protected() {
				p.Init(r)
			}
		}
	})
	//init plugins
	for _, p := range e.Plugins {
		if !p.Protected() {
			p.Init(r)
		}
	}
	l.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				fmt.Printf("Now server is running on %s\n", c.HostPort)
				fmt.Printf("Test with Get: curl -g 'http://%s/api/v1/graphql?query={hello}'\n", c.HostPort)
				_ = http.ListenAndServe(c.HostPort, r)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			db := s.GetDB()
			if db != nil {
				_ = db.Close()
			}
			fmt.Printf("%s shutdown complete", c.AppName)
			return nil
		},
	})
}
