package token

import (
	"context"
	"sync"
	"time"

	"github.com/Lumiaqian/go-sdk-core/cache"
)

type DefaultTokenProvider struct {
	mutex   sync.RWMutex
	cache   cache.Cache
	fetcher TokenFetcher
}

func NewDefaultTokenProvider(cache cache.Cache, fetcher TokenFetcher) *DefaultTokenProvider {
	return &DefaultTokenProvider{
		cache:   cache,
		fetcher: fetcher,
	}
}

func (p *DefaultTokenProvider) GetAccessToken(ctx context.Context) (string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	key := p.fetcher.GenerateCacheKey()

	// Try to get token from cache
	token, err := p.cache.Get(ctx, key)
	if err == nil {
		return token, nil
	}

	// Fetch token from server
	token, expiry, err := p.fetcher.FetchToken(ctx)
	if err != nil {
		return "", err
	}

	// Save token to cache
	err = p.cache.Set(ctx, key, token, time.Duration(expiry)*time.Second)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (p *DefaultTokenProvider) RefreshAccessToken(ctx context.Context) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	key := p.fetcher.GenerateCacheKey()
	token, expiry, err := p.fetcher.FetchToken(ctx)
	if err != nil {
		return "", err
	}

	err = p.cache.Set(ctx, key, token, time.Duration(expiry)*time.Second)
	if err != nil {
		return "", err
	}

	return token, nil
}
