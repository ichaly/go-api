package plugin

import (
	"github.com/dosco/graphjin/serv"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

type CaptchaService struct {
	Router  *chi.Mux
	Config  *base.Config
	Service *serv.Service
}

func NewCaptchaService(c *base.Config, r *chi.Mux, s *serv.Service) core.Plugin {
	return &CaptchaService{
		Router:  r,
		Config:  c,
		Service: s,
	}
}

func (my *CaptchaService) Name() string {
	return "Captcha Service"
}

func (my *CaptchaService) Init() {
	driver := my.Config.Captcha
	captcha := base64Captcha.NewCaptcha(driver, store)

	my.Router.Get("/captcha", func(w http.ResponseWriter, r *http.Request) {
		id, data, err := captcha.Generate()
		if err != nil {
			return
		}
		render.JSON(w, r, map[string]string{"id": id, "data": data})
	})
}
