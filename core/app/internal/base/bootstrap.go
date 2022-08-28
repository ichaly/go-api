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
	Plugins     []core.Plugin     `group:"plugin"`
	Middlewares []core.Middleware `group:"middleware"`
}

func Bootstrap(
	l fx.Lifecycle, s *serv.Service, r *chi.Mux, c *Config, e Enhance,
) {
	//init Plugins and Middlewares
	for _, p := range e.Plugins {
		p.Init()
	}
	for _, m := range e.Middlewares {
		m.Init()
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
