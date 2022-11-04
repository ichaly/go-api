package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/ichaly/go-api/core/app/pkg"
	"github.com/ichaly/go-api/core/app/pkg/render"
	"log"
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
		return render.JSON(w, core.OK.WithData(data))
	})
	o.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	o.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	return o
}
