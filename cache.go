package cache

import (
	"errors"

	"github.com/rahulbharuka/lru-ttl-cache/common"
	lruttl "github.com/rahulbharuka/lru-ttl-cache/lru_ttl"
)

// New returns a new in-memory cache of given size and cache type.
func New(size int, cacheType int) (common.Cache, error) {
	if size <= 0 {
		return nil, errors.New("invalid cache size")
	}

	switch cacheType {
	case common.CTLRUTTL:
		return lruttl.NewLRUTTLCache(size)
	}

	return nil, errors.New("invalid cache type")
}
