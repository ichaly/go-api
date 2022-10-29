package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"log"
)

func NewOauthServer(t oauth2.TokenStore, s oauth2.ClientStore) *server.Server {
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

	return o
}
