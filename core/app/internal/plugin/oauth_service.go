package plugin

import (
	"github.com/eko/gocache/v3/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/wenlng/go-captcha/captcha"
	"net/http"
)

type OauthService struct {
	Router  *chi.Mux
	Config  *base.Config
	Captcha *captcha.Captcha
	Store   *cache.Cache[string]
}

func NewOauthService(c *base.Config, r *chi.Mux, s *cache.Cache[string]) core.Plugin {
	g := captcha.GetCaptcha()
	return &OauthService{Store: s, Config: c, Router: r, Captcha: g}
}

func (my *OauthService) Name() string {
	return "OauthService"
}

func (my *OauthService) Init() {
	my.Router.Route("/oauth", func(r chi.Router) {
		r.Get("/authorize", my.authorizeHandler())
		r.Get("/token", my.tokenHandler())
	})
}

func (my *OauthService) authorizeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, base.OK)
	}
}

func (my *OauthService) tokenHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, base.OK)
	}
}
