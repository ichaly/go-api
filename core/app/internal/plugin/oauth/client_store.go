package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/store"
)

func NewOauthClientStore() oauth2.ClientStore {
	return store.NewClientStore()
}
