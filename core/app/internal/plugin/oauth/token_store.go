package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/store"
)

func NewOauthTokenStore() (oauth2.TokenStore, error) {
	return store.NewMemoryTokenStore()
}
