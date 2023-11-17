package token

import (
	"context"
	"go-sdk-core/cache"
	"sync"
	"time"
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
	// 从cache中获取
	token, err := p.cache.Get(ctx, key)
	if err == nil {
		return token, nil
	}
	// 从服务器中获取
	token, expiry, err := p.fetcher.FetchToken(ctx)
	if err != nil {
		return "", err
	}
	// 保存到cache中
	err = p.cache.Set(ctx, key, token, time.Duration(expiry)*time.Second)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (p *DefaultTokenProvider) RefreshAccessToken(ctx context.Context) (string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	key := p.fetcher.GenerateCacheKey()
	// 从服务器中获取
	token, expiry, err := p.fetcher.FetchToken(ctx)
	if err != nil {
		return "", err
	}
	// 保存到cache中
	err = p.cache.Set(ctx, key, token, time.Duration(expiry)*time.Second)
	if err != nil {
		return "", err
	}
	return token, nil
}
