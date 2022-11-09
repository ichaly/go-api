package oauth

import (
	"context"
	"fmt"
	"github.com/eko/gocache/v3/cache"
	cacheStore "github.com/eko/gocache/v3/store"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	"time"
)

const CachePrefix = "token"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type TokenStore struct {
	Cache       *cache.Cache[string]
	keyGenerate func(key string) string
}

func NewOauthTokenStore(c *cache.Cache[string]) oauth2.TokenStore {
	return &TokenStore{Cache: c, keyGenerate: func(key string) string {
		return fmt.Sprintf("%s:%s", CachePrefix, key)
	}}
}

func (my *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	jv, err := json.MarshalToString(info)
	if err != nil {
		return err
	}

	if code := info.GetCode(); code != "" {
		return my.Cache.Set(ctx, my.keyGenerate(code), jv, cacheStore.WithExpiration(info.GetCodeExpiresIn()))
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
			if err := my.Cache.Set(ctx, my.keyGenerate(refresh), basicID, cacheStore.WithExpiration(rexp)); err != nil {
				return err
			}
		}
	}

	if err = my.Cache.Set(ctx, my.keyGenerate(basicID), jv, cacheStore.WithExpiration(rexp)); err != nil {
		return err
	}

	return my.Cache.Set(ctx, my.keyGenerate(info.GetAccess()), basicID, cacheStore.WithExpiration(aexp))
}

// RemoveByCode delete the authorization code
func (my *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	return my.Cache.Delete(ctx, my.keyGenerate(code))
}

// RemoveByAccess use the access token to delete the token information
func (my *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return my.Cache.Delete(ctx, my.keyGenerate(access))
}

// RemoveByRefresh use the refresh token to delete the token information
func (my *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return my.Cache.Delete(ctx, my.keyGenerate(refresh))
}

// GetByCode use the authorization code for token information data
func (my *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return my.getData(ctx, code)
}

// GetByAccess use the access token for token information data
func (my *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	basicID, err := my.getBasicID(ctx, access)
	if err != nil {
		return nil, err
	}
	return my.getData(ctx, basicID)
}

// GetByRefresh use the refresh token for token information data
func (my *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	basicID, err := my.getBasicID(ctx, refresh)
	if err != nil {
		return nil, err
	}
	return my.getData(ctx, basicID)
}

func (my *TokenStore) getData(ctx context.Context, key string) (oauth2.TokenInfo, error) {
	val, err := my.Cache.Get(ctx, my.keyGenerate(key))
	if err != nil {
		return nil, err
	}
	var t *models.Token
	err = json.UnmarshalFromString(val, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (my *TokenStore) getBasicID(ctx context.Context, key string) (string, error) {
	return my.Cache.Get(ctx, my.keyGenerate(key))
}
