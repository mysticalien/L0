package tests_test

import (
	"testing"

	"L0/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	c := cache.NewCache()

	orderID := "12345"
	orderData := []byte(`{"order_uid": "12345", "track_number": "TRACK123"}`)

	c.Set(orderID, orderData)

	retrievedData, exists := c.Get(orderID)
	assert.True(t, exists)
	assert.Equal(t, orderData, retrievedData)
}

func TestCache_GetNonExisting(t *testing.T) {
	c := cache.NewCache()

	_, exists := c.Get("non-existing-id")
	assert.False(t, exists)
}
