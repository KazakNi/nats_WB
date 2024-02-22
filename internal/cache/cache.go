package cache

import (
	"log/slog"
	"sync"
)

type Cache struct {
	sync.RWMutex
	items map[string][]byte
}

func New() *Cache {
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

}

func (c *Cache) Warm() {
	// get slice of IDs
	// iterate over IDs and get all data
}
