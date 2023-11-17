package token

import (
	"context"
)

type TokenProvider interface {
	// GetAccessToken 从缓存或远程服务器获取access token。
	GetAccessToken(ctx context.Context) (string, error)

	// RefreshAccessToken 强制刷新access token，无论当前token是否已过期。
	RefreshAccessToken(ctx context.Context) (string, error)
}

type TokenFetcher interface {
	// FetchToken 从远程服务器获取access token
	FetchToken(ctx context.Context) (token string, expiry int64, err error)
	// GenerateCacheKey 生成缓存key
	GenerateCacheKey() string
}
