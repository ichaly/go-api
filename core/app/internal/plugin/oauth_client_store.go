package plugin

import (
	"github.com/go-oauth2/oauth2/v4/store"
)

func NewOauthClientStore() *store.ClientStore {
	return store.NewClientStore()
}
