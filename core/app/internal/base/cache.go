package base

import (
	"github.com/allegro/bigcache"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewCache(c *Config) (*cache.Cache[string], error) {
	var s store.StoreInterface
	if c.Cache.Type == "redis" {
		s = store.NewRedis(redis.NewClient(&redis.Options{
			Addr:     c.Cache.Url,
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
