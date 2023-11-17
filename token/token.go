package token

import (
	"context"
)

type TokenProvider interface {
	// GetAccessToken retrieves the access token from either the cache or the remote server.
	GetAccessToken(ctx context.Context) (string, error)

	// RefreshAccessToken forcefully refreshes the access token, regardless of its expiration status.
	RefreshAccessToken(ctx context.Context) (string, error)
}

type TokenFetcher interface {
	// FetchToken retrieves an access token from a remote server
	FetchToken(ctx context.Context) (token string, expiry int64, err error)
	// GenerateCacheKey generates a cache key
	GenerateCacheKey() string
}
