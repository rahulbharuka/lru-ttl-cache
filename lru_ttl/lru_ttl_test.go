package lruttl

import (
	"testing"
	"time"

	"github.com/rahulbharuka/lru-ttl-cache/common"
	"github.com/stretchr/testify/assert"
)

func TestNewLRUTTLCache(t *testing.T) {
	lruCache, err := NewLRUTTLCache(2)
	_ = assert.NoError(t, err) &&
		assert.NotNil(t, lruCache)

	t.Run("empty-cache", func(t *testing.T) {
		_, err := lruCache.Get("sg")
		assert.Error(t, err)
	})

	t.Run("happy-path", func(t *testing.T) {
		err := lruCache.Set("sg", "singapore", common.NoExpiry)
		assert.NoError(t, err)

		err = lruCache.Set("id", "indonesia", 1*time.Second)
		assert.NoError(t, err)

		val, err := lruCache.Get("id")
		_ = assert.NoError(t, err) &&
			assert.Equal(t, "indonesia", val)
	})

	t.Run("ttl-expiry", func(t *testing.T) {
		time.Sleep(1 * time.Second)
		lruCache.Set("ms", "malaysia", common.NoExpiry)
		val, err := lruCache.Get("sg")
		_ = assert.NoError(t, err) &&
			assert.Equal(t, "singapore", val)
	})
}
