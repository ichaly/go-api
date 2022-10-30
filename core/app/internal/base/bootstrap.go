package base

import (
	"context"
	"fmt"
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi/v5"
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
		m.Init()
	}
	l.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			//init plugins
			for _, p := range e.Plugins {
				p.Init()
			}
			//attach graphql
			r.Group(func(r chi.Router) {
				r.Use(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						ctx := context.WithValue(r.Context(), "user", "123")
						//if _, err := my.Oauth.ValidationBearerToken(r); err != nil {
						//	render.JSON(w, r, base.ERROR.WithData(err.Error()))
						//	return
						//}
						next.ServeHTTP(w, r.WithContext(ctx))
					})
				})
				_ = s.Attach(r)
			})
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
