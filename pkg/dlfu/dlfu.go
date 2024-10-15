package dlfu

import (
	"context"
	"sort"
	"time"

	"github.com/puzpuzpuz/xsync/v3"

	"github.com/Aki0x137/wind/config"
	"github.com/Aki0x137/wind/types"
)

type Item[V any] struct {
	key    string
	value  V
	score  types.Float64
	expiry time.Time
}

func (i *Item[V]) isExpired() bool {
	return time.Now().After(i.expiry)
}

func (i *Item[V]) Value() V {
	return i.value
}

type DLFUCache[V any] struct {
	capacity      int
	weight        float64
	decay         float64
	increment     types.Float64
	expiryEnabled bool

	cache *xsync.MapOf[string, *Item[V]]
}

func NewDLFUCache[V any](config config.DLFUConfig) *DLFUCache[V] {
	cache := &DLFUCache[V]{
		capacity:      config.Capacity,
		weight:        config.Weight,
		expiryEnabled: config.ExpiryEnabled,
		cache:         xsync.NewMapOf[string, *Item[V]](),
	}

	cache.increment.Store(1.0) // Initial value

	if config.Weight == 0.0 { // behaves like LFU cache
		cache.decay = 1.0
	} else {
		p := float64(config.Capacity) * config.Weight
		cache.decay = (p + 1.0) / p
	}

	// TODO: set timer for trimmer call using config.TrimInterval

	return cache
}

func (c *DLFUCache[V]) Set(items map[string]V, expiry time.Duration) {

	expiresAt := time.Now().Add(expiry)
	for key, val := range items {
		if item, ok := c.cache.Load(key); ok {
			item.value = val
			item.expiry = time.Now().Add(expiry)
			continue
		}
		item := &Item[V]{
			key:    key,
			value:  val,
			score:  c.increment,
			expiry: expiresAt,
		}
		c.cache.Store(key, item)
	}
}

// Get returns map of keys for which values are found in Cache and slice of keys for which value is not found in cache
func (c *DLFUCache[V]) Get(ctx context.Context, keys []string) (map[string]V, []string) {
	result := make(map[string]V)
	missingKeys := make([]string, 0)
	increment := c.increment.Load()

	for i, key := range keys {
		if ctx.Err() != nil {
			return result, append(keys[i:], missingKeys...)
		}
		if item, ok := c.cache.Load(key); ok && !item.isExpired() {
			result[key] = item.value
			item.score.Add(increment)
		} else {
			missingKeys = append(missingKeys, key)
		}
	}
	c.increment.Store(increment * c.decay)

	return result, missingKeys
}

func (c *DLFUCache[V]) trimmer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(250 * time.Millisecond):
			if ctx.Err() != nil {
				return
			}
			c.trim()
		}
	}
}

func (c *DLFUCache[V]) trim() {
	size := c.cache.Size()
	if size <= c.capacity {
		return
	}

	items := make([]*Item[V], 0, size)
	c.cache.Range(func(key string, value *Item[V]) bool {
		items = append(items, value)
		return true
	})
	sort.Sort(items)

	for i := 0; i < len(items)-c.capacity; i++ {
		c.cache.Delete(items[i].key)
	}
}
