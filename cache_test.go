package cache

import (
	"testing"
	"time"

	"github.com/rahulbharuka/lru-ttl-cache/common"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	t.Run("invalid-size", func(t *testing.T) {
		lruCache, err := New(-1, common.CTLRUTTL)
		_ = assert.Error(t, err) &&
			assert.Nil(t, lruCache)

	})

	t.Run("invalid-cache-type", func(t *testing.T) {
		lruCache, err := New(10, common.CTUndefined)
		_ = assert.Error(t, err) &&
			assert.Nil(t, lruCache)

	})

	t.Run("happy-path", func(t *testing.T) {
		lruCache, err := New(2, common.CTLRUTTL)
		_ = assert.NoError(t, err) &&
			assert.NotNil(t, lruCache)

		err = lruCache.Set("sg", "singapore", common.NoExpiry)
		assert.NoError(t, err)

		err = lruCache.Set("id", "indonesia", 1*time.Second)
		assert.NoError(t, err)

		val, err := lruCache.Get("sg")
		_ = assert.NoError(t, err) &&
			assert.Equal(t, "singapore", val)
	})
}
