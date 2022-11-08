package oauth

import (
	"context"
	"github.com/eko/gocache/v3/cache"
	cacheStore "github.com/eko/gocache/v3/store"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type TokenStore struct {
	Cache *cache.Cache[string]
}

func NewOauthTokenStore(c *cache.Cache[string]) oauth2.TokenStore {
	if tokenStore, err := store.NewMemoryTokenStore(); err == nil {
		return tokenStore
	}
	return &TokenStore{Cache: c}
}

func (my *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	jv, err := json.MarshalToString(info)
	if err != nil {
		return err
	}

	if code := info.GetCode(); code != "" {
		return my.Cache.Set(ctx, code, jv, cacheStore.WithExpiration(info.GetCodeExpiresIn()))
	}

	basicID := uuid.Must(uuid.NewRandom()).String()
	aexp := info.GetAccessExpiresIn()
	rexp := aexp

	if refresh := info.GetRefresh(); refresh != "" {
		ct := time.Now()
		rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
		if aexp.Seconds() > rexp.Seconds() {
			aexp = rexp
		}
		if info.GetRefreshExpiresIn() != 0 {
			if err := my.Cache.Set(ctx, refresh, basicID, cacheStore.WithExpiration(rexp)); err != nil {
				return err
			}
		}
	}

	if err = my.Cache.Set(ctx, basicID, jv, cacheStore.WithExpiration(rexp)); err != nil {
		return err
	}

	return my.Cache.Set(ctx, info.GetAccess(), basicID, cacheStore.WithExpiration(aexp))
}

// RemoveByCode delete the authorization code
func (my *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	return my.Cache.Delete(ctx, code)
}

// RemoveByAccess use the access token to delete the token information
func (my *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return my.Cache.Delete(ctx, access)
}

// RemoveByRefresh use the refresh token to delete the token information
func (my *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return my.Cache.Delete(ctx, refresh)
}

// GetByCode use the authorization code for token information data
func (my *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// GetByAccess use the access token for token information data
func (my *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// GetByRefresh use the refresh token for token information data
func (my *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	return nil, nil
}
