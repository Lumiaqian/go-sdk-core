package cache

import (
	"context"
	"errors"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Memcache struct {
	cache *gocache.Cache
}

func NewMemcache(defaultExpiration, cleanupInterval time.Duration) *Memcache {
	return &Memcache{
		cache: gocache.New(defaultExpiration, cleanupInterval),
	}
}

func (g *Memcache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		g.cache.Set(key, value, expiration)
		return nil
	}
}

func (g *Memcache) Get(ctx context.Context, key string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		val, found := g.cache.Get(key)
		if !found {
			return "", errors.New("key not found")
		}
		return val.(string), nil
	}
}

func (g *Memcache) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		g.cache.Delete(key)
		return nil
	}
}
