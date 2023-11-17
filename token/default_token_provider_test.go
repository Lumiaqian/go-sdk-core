package token

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCache is a mock implementation of the Cache interface
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// MockTokenFetcher is a mock implementation of the TokenFetcher interface
type MockTokenFetcher struct {
	mock.Mock
}

func (m *MockTokenFetcher) FetchToken(ctx context.Context) (string, int64, error) {
	args := m.Called(ctx)
	return args.String(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockTokenFetcher) GenerateCacheKey() string {
	args := m.Called()
	return args.String(0)
}

func TestNewDefaultTokenProvider(t *testing.T) {
	mockCache := new(MockCache)
	mockFetcher := new(MockTokenFetcher)

	provider := NewDefaultTokenProvider(mockCache, mockFetcher)
	assert.NotNil(t, provider, "NewDefaultTokenProvider returned nil")
	assert.Equal(t, mockCache, provider.cache, "NewDefaultTokenProvider did not correctly assign cache field")
	assert.Equal(t, mockFetcher, provider.fetcher, "NewDefaultTokenProvider did not correctly assign fetcher field")
}

// Additional tests for GetAccessToken and RefreshAccessToken can be added here

func TestDefaultTokenProvider_GetAccessToken_CacheHit(t *testing.T) {
	mockCache := new(MockCache)
	mockFetcher := new(MockTokenFetcher)
	provider := NewDefaultTokenProvider(mockCache, mockFetcher)

	// Set the expected behavior for the GenerateCacheKey method.
	mockFetcher.On("GenerateCacheKey").Return("mockedCacheKey")

	// Set the expected behavior for the Get method.
	mockCache.On("Get", mock.Anything, "mockedCacheKey").Return("cachedToken", nil)

	token, err := provider.GetAccessToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "cachedToken", token)
}

func TestDefaultTokenProvider_GetAccessToken_FetcherSuccess(t *testing.T) {
	mockCache := new(MockCache)
	mockFetcher := new(MockTokenFetcher)
	provider := NewDefaultTokenProvider(mockCache, mockFetcher)

	// Set the expected behavior for the GenerateCacheKey method.
	mockFetcher.On("GenerateCacheKey").Return("mockedCacheKey")
	mockCache.On("Get", mock.Anything, "mockedCacheKey").Return("", errors.New("cache miss"))
	mockFetcher.On("FetchToken", mock.Anything).Return("fetchedToken", int64(3600), nil)
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	token, err := provider.GetAccessToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "fetchedToken", token)
}

func TestDefaultTokenProvider_GetAccessToken_Failure(t *testing.T) {
	mockCache := new(MockCache)
	mockFetcher := new(MockTokenFetcher)
	provider := NewDefaultTokenProvider(mockCache, mockFetcher)

	// Set the expected behavior for the GenerateCacheKey method.
	mockFetcher.On("GenerateCacheKey").Return("mockedCacheKey")
	mockCache.On("Get", mock.Anything, "mockedCacheKey").Return("", errors.New("cache miss"))
	mockFetcher.On("FetchToken", mock.Anything).Return("", int64(0), errors.New("fetcher error"))

	token, err := provider.GetAccessToken(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "", token)
}

func TestDefaultTokenProvider_RefreshAccessToken_Success(t *testing.T) {
	mockCache := new(MockCache)
	mockFetcher := new(MockTokenFetcher)
	provider := NewDefaultTokenProvider(mockCache, mockFetcher)

	mockFetcher.On("FetchToken", mock.Anything).Return("refreshedToken", int64(3600), nil)
	// Set the expected behavior for the GenerateCacheKey method.
	mockFetcher.On("GenerateCacheKey").Return("mockedCacheKey")
	mockCache.On("Set", mock.Anything, "mockedCacheKey", mock.Anything, mock.Anything).Return(nil)

	token, err := provider.RefreshAccessToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "refreshedToken", token)
}

func TestDefaultTokenProvider_RefreshAccessToken_Failure(t *testing.T) {
	mockCache := new(MockCache)
	mockFetcher := new(MockTokenFetcher)
	provider := NewDefaultTokenProvider(mockCache, mockFetcher)
	// Set the expected behavior for the GenerateCacheKey method.
	mockFetcher.On("GenerateCacheKey").Return("mockedCacheKey")
	mockFetcher.On("FetchToken", mock.Anything).Return("", int64(0), errors.New("fetcher error"))

	token, err := provider.RefreshAccessToken(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "", token)
}
