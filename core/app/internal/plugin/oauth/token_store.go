package oauth

import (
	"github.com/eko/gocache/v3/cache"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/store"
)

func NewOauthTokenStore(s *cache.Cache[string]) (oauth2.TokenStore, error) {
	return store.NewMemoryTokenStore()
}
