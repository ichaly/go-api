package oauth

import (
	"github.com/eko/gocache/v3/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/ichaly/go-api/core/app/internal/base"
	"github.com/ichaly/go-api/core/app/pkg"
	"log"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
)

type OauthService struct {
	Router *chi.Mux
	Oauth  *server.Server
	Store  *cache.Cache[string]
}

func NewOauthService(r *chi.Mux, t oauth2.TokenStore, s oauth2.ClientStore) core.Plugin {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(t, nil)
	manager.MapClientStorage(s)

	o := server.NewDefaultServer(manager)
	o.SetAllowGetAccessRequest(true)
	o.SetClientInfoHandler(server.ClientFormHandler)
	o.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	o.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	return &OauthService{Oauth: o, Router: r}
}

func (my *OauthService) Name() string {
	return "OauthService"
}

func (my *OauthService) Init() {
	//使用中间件过滤
	//my.Router.Use(func(next http.Handler) http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		ctx := context.WithValue(r.Context(), "user", "123")
	//		if _, err := my.Oauth.ValidationBearerToken(r); err != nil {
	//			render.JSON(w, r, base.ERROR.WithData(err.Error()))
	//			return
	//		}
	//		next.ServeHTTP(w, r.WithContext(ctx))
	//	})
	//})
	my.Router.Group(func(r chi.Router) {
		r.Route("/oauth", func(r chi.Router) {
			r.Get("/token", my.tokenHandler())
			r.Get("/authorize", my.authorizeHandler())
		})
	})
}

func (my *OauthService) authorizeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := my.Oauth.HandleAuthorizeRequest(w, r); err != nil {
			render.JSON(w, r, base.ERROR.WithData(err.Error()))
		}
	}
}

func (my *OauthService) tokenHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := my.Oauth.HandleTokenRequest(w, r); err != nil {
			render.JSON(w, r, base.ERROR.WithData(err.Error()))
		}
	}
}
