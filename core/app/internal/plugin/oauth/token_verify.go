package oauth

import (
	"context"
	gql "github.com/dosco/graphjin/core"
	"github.com/go-chi/chi"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/ichaly/go-api/core/app/pkg/render"
	"net/http"
)

type verify struct {
	Oauth *server.Server
}

func NewOauthVerify(s *server.Server) core.Plugin {
	return &verify{Oauth: s}
}

func (my *verify) Name() string {
	return "OauthVerify"
}

func (my *verify) Protected() bool {
	return true
}

func (my *verify) Init(r chi.Router) {
	//使用中间件鉴权
	r.Use(my.verifyHandler)
}

func (my *verify) verifyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, err := my.Oauth.ValidationBearerToken(r); err != nil {
			_ = render.JSON(w, core.ERROR.AddError(err.Error()))
		} else {
			ctx := context.WithValue(r.Context(), gql.UserIDKey, token.GetUserID())
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
