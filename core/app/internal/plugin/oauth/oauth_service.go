package oauth

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/pkg"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/server"
)

type OauthService struct {
	Router *chi.Mux
	Oauth  *server.Server
}

func NewOauthService(r *chi.Mux, o *server.Server) core.Plugin {
	return &OauthService{Oauth: o, Router: r}
}

func (my *OauthService) Name() string {
	return "OauthService"
}

func (my *OauthService) Init() {
	//授权路由
	my.Router.Route("/oauth", func(r chi.Router) {
		r.HandleFunc("/token", my.tokenHandler())
		r.HandleFunc("/authorize", my.authorizeHandler())
	})
}

func (my *OauthService) tokenHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := my.Oauth.HandleTokenRequest(w, r); err != nil {
			render.JSON(w, r, base.ERROR.WithData(err.Error()))
		}
	}
}

func (my *OauthService) authorizeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := my.Oauth.HandleAuthorizeRequest(w, r); err != nil {
			render.JSON(w, r, base.ERROR.WithData(err.Error()))
		}
	}
}
