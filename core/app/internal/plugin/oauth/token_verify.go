package oauth

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/ichaly/go-api/core/app/internal/base"
	core "github.com/ichaly/go-api/core/app/pkg"
	"net/http"
)

type TokenVerify struct {
	Router *chi.Mux
	Oauth  *server.Server
}

func NewOauthTokenVerify(r *chi.Mux, s *server.Server) core.Plugin {
	return &TokenVerify{Router: r, Oauth: s}
}

func (my *TokenVerify) Name() string {
	return "OauthTokenVerify"
}

func (my *TokenVerify) Init() {
	//使用中间件鉴权
	//my.Router.Use(my.verifyHandler)
	//my.Router.Route("/", func(r chi.Router) {
	//	r.Use(my.verifyHandler)
	//})
}

func (my *TokenVerify) verifyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "user", "123")
		if _, err := my.Oauth.ValidationBearerToken(r); err != nil {
			render.JSON(w, r, base.ERROR.WithData(err.Error()))
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
