package cache

import (
	"L0/internal/storage"
	"context"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	cache map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string][]byte),
	}
}

func (c *Cache) Set(orderID string, orderData []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[orderID] = orderData
}

func (c *Cache) Get(orderID string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, exists := c.cache[orderID]
	return data, exists
}

func (c *Cache) LoadFromDB(orders map[string][]byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, data := range orders {
		c.cache[id] = data
	}
}

func (c *Cache) RestoreFromDB(db *storage.Storage, ctx context.Context) error {
	orders, err := db.GetAll(ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		c.cache[order.OrderUID] = order.OrderInfo
	}
	return nil
}
