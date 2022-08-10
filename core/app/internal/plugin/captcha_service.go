package plugin

import (
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/pkg"
	"net/http"
)

func NewCaptchaService(c *base.Config, r *chi.Mux, s *serv.Service) core.Plugin {
	return &CaptchaService{
		Router:  r,
		Config:  c,
		Service: s,
	}
}

type CaptchaService struct {
	Router  *chi.Mux
	Config  *base.Config
	Service *serv.Service
}

func (my *CaptchaService) Name() string {
	return "Captcha Service"
}

func (my *CaptchaService) Init() {
	my.Router.Get("/captcha", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "successfully created"})
	})
}
