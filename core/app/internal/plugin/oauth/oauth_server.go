package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/ichaly/go-api/core/app/pkg/render"
	"net/http"
)

func NewOauthServer(t oauth2.TokenStore, s oauth2.ClientStore) *server.Server {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(t, nil)
	manager.MapClientStorage(s)

	o := server.NewDefaultServer(manager)
	o.SetAllowGetAccessRequest(true)
	o.SetClientInfoHandler(server.ClientFormHandler)
	o.SetResponseTokenHandler(func(w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) error {
		code := 200
		if len(statusCode) > 0 {
			code = statusCode[0]
		}
		var err interface{}
		if v, e := data["error"]; e {
			err = v
		}
		if code == 200 {
			return render.JSON(w, core.OK.WithData(data))
		} else {
			return render.JSON(w, core.NewResult(code).AddError(err))
		}
	})
	o.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		return errors.NewResponse(err, http.StatusInternalServerError)
	})

	return o
}
