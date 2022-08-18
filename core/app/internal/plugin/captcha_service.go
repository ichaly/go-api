package plugin

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

type CaptchaService struct {
	Router  *chi.Mux
	Config  *base.Config
	Captcha *base64Captcha.Captcha
}

func NewCaptchaService(c *base.Config, r *chi.Mux) core.Plugin {
	b := base64Captcha.NewCaptcha(
		c.Driver, base64Captcha.DefaultMemStore,
	)
	return &CaptchaService{Config: c, Router: r, Captcha: b}
}

func (my *CaptchaService) Name() string {
	return "Captcha Service"
}

func (my *CaptchaService) Init() {
	my.Router.Get("/captcha", func(w http.ResponseWriter, r *http.Request) {
		id, data, err := my.Captcha.Generate()
		if err != nil {
			return
		}
		render.JSON(w, r, map[string]string{"id": id, "data": data})
	})
}
