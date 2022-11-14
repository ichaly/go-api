package base

import (
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

func NewStore(c *Config) (*cache.Cache[string], error) {
	var s store.StoreInterface
	if strings.ToLower(c.Cache.Type) == "redis" {
		args := []interface{}{c.Cache.Host, c.Cache.Port}
		s = store.NewRedis(redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", args...),
			Username: c.Cache.Username,
			Password: c.Cache.Password,
		}))
	} else {
		client, err := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
		if err != nil {
			return nil, err
		}
		s = store.NewBigcache(client)
	}
	return cache.New[string](s), nil
}
