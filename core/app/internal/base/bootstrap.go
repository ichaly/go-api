package base

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/json-iterator/go/extra"
	"go.uber.org/fx"
	"net/http"
)

func init() {
	extra.RegisterFuzzyDecoders()
}

type Enhance struct {
	fx.In
	Plugins     []core.Plugin `group:"plugin"`
	Middlewares []core.Plugin `group:"middleware"`
}

func Bootstrap(
	l fx.Lifecycle, c *Config, e Enhance, r *chi.Mux, s *Engine,
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
		s.Attach(r)
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
			fmt.Printf("%s shutdown complete", c.AppName)
			return nil
		},
	})
}
