package cache

import (
	"log/slog"
	"sync"
)

type Cache struct {
	sync.RWMutex
	Items map[string][]byte
}

func New() *Cache {
	items := make(map[string][]byte)

	cache := Cache{
		Items: items,
	}

	//start warming

	return &cache
}

func (c *Cache) Get(key string) (data []byte, found bool) {
	c.RLock()

	defer c.RUnlock()

	value, found := c.Items[key]

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

	c.Items[key] = data
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
