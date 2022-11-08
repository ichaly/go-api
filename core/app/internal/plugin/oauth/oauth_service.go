package oauth

import (
	"github.com/go-chi/chi"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/ichaly/go-api/core/app/pkg/render"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/server"
)

type OauthService struct {
	Oauth *server.Server
}

func NewOauthService(o *server.Server) core.Plugin {
	return &OauthService{Oauth: o}
}

func (my *OauthService) Name() string {
	return "OauthService"
}

func (my *OauthService) Protected() bool {
	return false
}

func (my *OauthService) Init(r chi.Router) {
	//授权路由
	r.Route("/oauth", func(r chi.Router) {
		r.HandleFunc("/token", my.tokenHandler())
		r.HandleFunc("/authorize", my.authorizeHandler())
	})
}

func (my *OauthService) tokenHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := my.Oauth.HandleTokenRequest(w, r); err != nil {
			_ = render.JSON(w, core.ERROR.AddError(err.Error()))
		}
	}
}

func (my *OauthService) authorizeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := my.Oauth.HandleAuthorizeRequest(w, r); err != nil {
			_ = render.JSON(w, core.ERROR.AddError(err.Error()))
		}
	}
}
