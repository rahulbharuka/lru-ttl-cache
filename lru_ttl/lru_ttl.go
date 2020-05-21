package lruttl

import (
	"container/heap"
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/rahulbharuka/lru-ttl-cache/common"
)

// lruEntry is the entry object for LRU-TTL cache.
type lruEntry struct {
	*common.Entry
	listPtr  *list.Element
	queueIdx int
}

// lruCache holds all data structures required for LRU-TTL cache.
type lruCache struct {
	size    int
	items   map[interface{}]*lruEntry
	lruList *list.List
	ttlQ    *priorityQueue
	lock    sync.RWMutex
}

// NewLRUTTLCache creates a new LRU-TTL cache of given size.
func NewLRUTTLCache(size int) (common.Cache, error) {
	pQueue := &priorityQueue{}
	heap.Init(pQueue)
	return &lruCache{
		size:    size,
		items:   map[interface{}]*lruEntry{},
		lruList: list.New(),
		ttlQ:    pQueue,
	}, nil
}

// addItem adds the item to all auxillary data structures of LRU-TTL cache.
func (c *lruCache) addItem(key interface{}, val interface{}, ttl time.Duration) error {
	item := &lruEntry{
		Entry: &common.Entry{
			Key: key,
			Val: val,
			TTL: ttl,
		},
	}

	if ttl != common.NoExpiry {
		item.ExpiryTime = time.Now().Add(ttl)
		heap.Push(c.ttlQ, item) // add only items with valid TTL to the queue
	}

	item.listPtr = c.lruList.PushFront(item) // mark the item as most recently used.
	c.items[key] = item
	return nil
}

// updateItem updates the item in all auxillary data structures of LRU-TTL cache.
func (c *lruCache) updateItem(item *lruEntry, val interface{}, ttl time.Duration) error {
	item.Val = val

	// remove item from TTL queue if it has been changed to no-expiry.
	if item.TTL != common.NoExpiry && ttl == common.NoExpiry {
		heap.Remove(c.ttlQ, item.queueIdx)
	} else if item.TTL == common.NoExpiry && ttl != common.NoExpiry {
		item.TTL = ttl
		heap.Push(c.ttlQ, item) // add item to the TTL queue
	} else if ttl != common.NoExpiry {
		item.TTL = ttl
		item.ExpiryTime = time.Now().Add(ttl)
		heap.Fix(c.ttlQ, item.queueIdx) // update TTL min heap
	}

	c.lruList.MoveToFront(item.listPtr) // update LRU linked list
	return nil
}

// removeItem removes the item from all auxillary data structures.
func (c *lruCache) removeItem(item *lruEntry) error {
	if item.TTL != common.NoExpiry {
		heap.Remove(c.ttlQ, item.queueIdx) // remove from TTL priority queue
	}
	c.lruList.Remove(item.listPtr) // remove from LRU linked list
	delete(c.items, item.Key)      // remove from items map
	return nil
}

// Get returns the value corresponding to given key from the cache if available.
func (c *lruCache) Get(key interface{}) (interface{}, error) {
	item, ok := c.items[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	// if the item has expired, remove it before returning an error
	if item.TTL > 0 && item.ExpiryTime.Before(time.Now()) {
		c.removeItem(item)
		return nil, errors.New("key not found")
	}

	c.lruList.MoveToFront(item.listPtr) // Mark the item as most recently used by moving it to front of the queue.
	// update lruList and ttlQueue
	return item.Val, nil
}

// Set adds/updates an item in the cache.
func (c *lruCache) Set(key interface{}, val interface{}, ttl time.Duration) error {
	// if the key exist, overwrite the item
	c.lock.Lock()
	defer c.lock.Unlock()

	if item, ok := c.items[key]; ok {
		return c.updateItem(item, val, ttl)
	}

	if len(c.items) < c.size {
		return c.addItem(key, val, ttl)
	}

	// cache is full; so check whether an item can be expired
	if len(*c.ttlQ) > 0 && (*c.ttlQ)[0].ExpiryTime.Before(time.Now()) {
		c.removeItem((*c.ttlQ)[0]) // delete the item corresponding to root of ttlQ
		return c.addItem(key, val, ttl)
	}

	// no item expired, so evict least recently used item from LRU linked list
	c.removeItem(c.lruList.Back().Value.(*lruEntry))
	return c.addItem(key, val, ttl)
}
