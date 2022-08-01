package internal

import (
	"context"
	"fmt"
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi"
	"go.uber.org/fx"
	"net/http"
)

func Bootstrap(lifecycle fx.Lifecycle, svc *serv.Service, mux *chi.Mux, cfg *serv.Config) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				fmt.Printf("Now server is running on %s\n", cfg.HostPort)
				fmt.Printf("Test with Get: curl -g 'http://%s/api/v1/graphql?query={hello}'\n", cfg.HostPort)
				_ = svc.Attach(mux)
				_ = http.ListenAndServe(cfg.HostPort, mux)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			db := svc.GetDB()
			if db != nil {
				_ = db.Close()
			}
			fmt.Printf("%s shutdown complete", cfg.AppName)
			return nil
		},
	})
}
