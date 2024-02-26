package cache

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"nats/api"
	"nats/internal/db"
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

	cache.Warm()

	return &cache
}

func (c *Cache) Get(key string) (data []byte, found bool) {
	c.RLock()

	defer c.RUnlock()

	value, found := c.Items[key]

	if !found {
		slog.Info(fmt.Sprintf("cannot found key:%s, checking database", key))
		var order api.Order
		order = db.GetItembyId(db.DBConnection, key)
		if len(order.Order_uid) == 0 {
			slog.Info(fmt.Sprintf("cannot found key:%s", key))
			return nil, false
		} else {
			order, err := json.Marshal(order)
			if err != nil {
				slog.Error(fmt.Sprintf("Error while marshalling in cache: %s", err))
			} else {
				return order, true
			}
		}
	}
	return value, true
}

func (c *Cache) Set(key string, data []byte) {
	c.Lock()

	defer c.Unlock()

	c.Items[key] = data
}

func (c *Cache) Warm() {

	items := db.GetItems(db.DBConnection)
	for _, val := range items {
		order, err := json.Marshal(val)
		if err != nil {
			slog.Error(fmt.Sprintf("Error while warming cache: %s", err))
		} else {
			c.Items[val.Order_uid] = order
		}
	}
}
