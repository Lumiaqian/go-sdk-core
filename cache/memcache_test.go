package cache

import (
	"context"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemcache(t *testing.T) {
	mc := NewMemcache(5*time.Minute, 10*time.Minute)
	require.NotNil(t, mc)
	assert.IsType(t, &Memcache{}, mc)
}

func TestMemcache_Set_Get_Delete(t *testing.T) {
	mc := NewMemcache(5*time.Minute, 10*time.Minute)
	ctx := context.Background()
	key := "testKey"
	value := "testValue"

	// Test Set
	err := mc.Set(ctx, key, value, cache.DefaultExpiration)
	assert.NoError(t, err)

	// Test Get
	val, err := mc.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, val)

	// Test Get with non-existing key
	_, err = mc.Get(ctx, "nonExistingKey")
	assert.Error(t, err)

	// Test Delete
	err = mc.Delete(ctx, key)
	assert.NoError(t, err)

	// Test Get after Delete
	_, err = mc.Get(ctx, key)
	assert.Error(t, err)
}

func TestMemcache_Set_Get_Delete_WithContextCancel(t *testing.T) {
	mc := NewMemcache(5*time.Minute, 10*time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	key := "testKey"
	value := "testValue"

	// Test Set with cancel
	cancel() // immediately cancel the context
	err := mc.Set(ctx, key, value, cache.DefaultExpiration)
	assert.Error(t, err)

	// Reset context
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Test Set
	err = mc.Set(ctx, key, value, cache.DefaultExpiration)
	assert.NoError(t, err)

	// Test Get with cancel
	cancel() // cancel the context before Get
	_, err = mc.Get(ctx, key)
	assert.Error(t, err)

	// Reset context
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Test Delete with cancel
	cancel() // cancel the context before Delete
	err = mc.Delete(ctx, key)
	assert.Error(t, err)
}
