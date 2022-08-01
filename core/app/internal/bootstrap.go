package internal

import (
	"context"
	"fmt"
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi"
	"github.com/ichaly/go-api/core/app/pkg"
	"go.uber.org/fx"
	"net/http"
)

type Params struct {
	fx.In
	Plugins     []core.Plugin     `group:"plugin"`
	Middlewares []core.Middleware `group:"middleware"`
}

func Bootstrap(
	lifecycle fx.Lifecycle, service *serv.Service, route *chi.Mux, config *serv.Config, params Params,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			//init Plugins and Middlewares
			for _, p := range params.Plugins {
				p.Init()
			}
			for _, m := range params.Middlewares {
				m.Init()
			}
			go func() {
				fmt.Printf("Now server is running on %s\n", config.HostPort)
				fmt.Printf("Test with Get: curl -g 'http://%s/api/v1/graphql?query={hello}'\n", config.HostPort)
				_ = http.ListenAndServe(config.HostPort, route)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			db := service.GetDB()
			if db != nil {
				_ = db.Close()
			}
			fmt.Printf("%s shutdown complete", config.AppName)
			return nil
		},
	})
}
