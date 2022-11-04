package oauth

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/ichaly/go-api/core/app/pkg/render"
	"net/http"
)

type TokenVerify struct {
	Oauth *server.Server
}

func NewOauthTokenVerify(s *server.Server) core.Plugin {
	return &TokenVerify{Oauth: s}
}

func (my *TokenVerify) Name() string {
	return "OauthTokenVerify"
}

func (my *TokenVerify) Protected() bool {
	return true
}

func (my *TokenVerify) Init(r chi.Router) {
	//使用中间件鉴权
	r.Use(my.verifyHandler)
}

func (my *TokenVerify) verifyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "user", "123")
		if _, err := my.Oauth.ValidationBearerToken(r); err != nil {
			_ = render.JSON(w, core.ERROR.WithData(err.Error()))
		} else {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
