package oauth

import (
	"context"
	"encoding/json"
	"github.com/eko/gocache/v3/cache"
	store2 "github.com/eko/gocache/v3/store"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/google/uuid"
	"time"
)

type CacheTokenStore struct {
	Store *cache.Cache[string]
}

func NewOauthTokenStore(s *cache.Cache[string]) (oauth2.TokenStore, error) {
	return store.NewMemoryTokenStore()
}

func (my *CacheTokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	jv, err := json.Marshal(info)
	if err != nil {
		return err
	}

	if code := info.GetCode(); code != "" {
		return my.Store.Set(ctx, code, string(jv), store2.WithExpiration(info.GetCodeExpiresIn()))
	}

	basicID := uuid.Must(uuid.NewRandom()).String()
	aexp := info.GetAccessExpiresIn()
	ct := time.Now()
	rexp := aexp

	if refresh := info.GetRefresh(); refresh != "" {
		rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
		if aexp.Seconds() > rexp.Seconds() {
			aexp = rexp
		}
		if err = my.Store.Set(ctx, refresh, basicID, store2.WithExpiration(rexp)); err != nil {
			return err
		}
	}

	if err = my.Store.Set(ctx, basicID, string(jv), store2.WithExpiration(rexp)); err != nil {
		return err
	}
	return my.Store.Set(ctx, info.GetAccess(), basicID, store2.WithExpiration(aexp))
}

// delete the authorization code
func (my *CacheTokenStore) RemoveByCode(ctx context.Context, code string) error {
	return nil
}

// use the access token to delete the token information
func (my *CacheTokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return nil
}

// use the refresh token to delete the token information
func (my *CacheTokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return nil
}

// use the authorization code for token information data
func (my *CacheTokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// use the access token for token information data
func (my *CacheTokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// use the refresh token for token information data
func (my *CacheTokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	return nil, nil
}
