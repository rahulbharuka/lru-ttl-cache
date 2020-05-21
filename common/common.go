package common

import "time"

const (
	// CTUndefined ...
	CTUndefined = iota
	// CTLRUTTL is a LRU-TTLcache
	CTLRUTTL
)

// NoExpiry means the cache entry will not expire.
const NoExpiry time.Duration = -1

// Cache is the interface for an in-memory cache.
type Cache interface {
	Get(key interface{}) (interface{}, error)
	Set(key interface{}, val interface{}, ttl time.Duration) error
}

// Entry object stores cache item details.
type Entry struct {
	Key        interface{}
	Val        interface{}
	TTL        time.Duration
	ExpiryTime time.Time
}
