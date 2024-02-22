package cache

import (
	"log/slog"
	"sync"
)

type Cache struct {
	sync.RWMutex
	items map[string][]byte
}

func New(c *Cache) *Cache {
	items := make(map[string][]byte)

	cache := Cache{
		items: items,
	}

	//start warming

	return &cache
}

func (c *Cache) Get(key string) (data []byte, found bool) {
	c.RLock()

	defer c.RUnlock()

	value, found := c.items[key]

	if !found {
		slog.Info("cannot found key:%s", key)
		return nil, false
	} else {
		return value, true
	}
}

func (c *Cache) Set(key string, data []byte) {
	c.Lock()

	defer c.Unlock()

	c.items[key] = data
}

func (c *Cache) Warm() {
	// db.getItems
	// process each item by Cache validation method
	// Set each validated item in bytes to the Cache
}

func (c *Cache) ValidateItem(rows string) []byte {
	// todo
	return nil
}
