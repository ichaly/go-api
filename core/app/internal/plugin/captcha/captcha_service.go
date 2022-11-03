package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eko/gocache/v3/cache"
	"github.com/go-chi/chi"
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/unrolled/render"
	"github.com/wenlng/go-captcha/captcha"
	"net/http"
	"strconv"
	"strings"
)

type CaptchaService struct {
	Config  *base.Config
	Render  *render.Render
	Captcha *captcha.Captcha
	Store   *cache.Cache[string]
}

func NewCaptchaService(c *base.Config, r *render.Render, s *cache.Cache[string]) core.Plugin {
	g := captcha.GetCaptcha()
	return &CaptchaService{Store: s, Config: c, Captcha: g, Render: r}
}

func (my *CaptchaService) Name() string {
	return "CaptchaService"
}

func (my *CaptchaService) Protected() bool {
	return false
}

func (my *CaptchaService) Init(r chi.Router) {
	r.Route("/captcha", func(r chi.Router) {
		r.Get("/verify", my.verifyHandler())
		r.Get("/generate", my.generateHandler())
	})
}

func (my *CaptchaService) Verify(ctx context.Context, key string, str string, del ...bool) (bool, error) {
	val, err := my.Store.Get(ctx, key)
	if err != nil {
		return false, err
	}

	if len(del) > 0 && del[0] {
		_ = my.Store.Delete(ctx, key)
	}

	var data map[int]captcha.CharDot
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return false, err
	}

	dots := strings.Split(str, ",")
	if (len(data) * 2) == len(dots) {
		return false, nil
	}

	for i, dot := range data {
		sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", dots[i*2]), 64)
		sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", dots[i*2+1]), 64)
		// 检测点位置,在原有的区域上添加额外边距进行扩张计算区域,不推荐设置过大的padding
		if !captcha.CheckPointDistWithPadding(
			int64(sx), int64(sy),
			int64(dot.Dx), int64(dot.Dy),
			int64(dot.Width), int64(dot.Height), 5,
		) {
			return false, nil
		}
	}

	return true, nil
}

func (my *CaptchaService) verifyHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		key := r.Form.Get("key")
		dots := r.Form.Get("dots")
		data, err := my.Verify(r.Context(), key, dots)
		if err != nil {
			panic(err)
		}
		_ = my.Render.JSON(w, http.StatusOK, core.OK.WithData(data))
	}
}

func (my *CaptchaService) generateHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dots, image, thumb, key, err := my.Captcha.Generate()
		if err != nil {
			panic(err)
		}
		raw, err := json.Marshal(dots)
		if err != nil {
			panic(err)
		}
		err = my.Store.Set(r.Context(), key, string(raw))
		if err != nil {
			panic(err)
		}
		_ = my.Render.JSON(w, http.StatusOK, core.OK.WithData(map[string]string{"key": key, "image": image, "thumb": thumb}))
	}
}
