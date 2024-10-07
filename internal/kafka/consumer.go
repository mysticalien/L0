package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/storage"
	"github.com/segmentio/kafka-go"
)

func saveWithRetry(ctx context.Context, db *storage.Storage, orderUID string, data []byte, maxRetries int, retryDelay time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = db.Save(ctx, orderUID, data)
		if err == nil {
			return nil
		}
		log.Printf("Retrying saving data... attempt %d", i+1)
		time.Sleep(retryDelay)
	}
	return err
}

func ConsumeOrders(cfg *config.KafkaConfig, cache *cache.Cache, dbStorage *storage.Storage) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Broker},
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message from Kafka: %v", err)
			continue
		}

		var orders []map[string]interface{}
		err = json.Unmarshal(msg.Value, &orders)
		if err == nil {
			log.Println("Processing JSON array of orders...")
			processOrders(orders, cache, dbStorage, cfg, msg.Value)
			continue
		}

		var order map[string]interface{}
		err = json.Unmarshal(msg.Value, &order)
		if err == nil {
			log.Println("Processing single order JSON...")
			processOrder(order, cache, dbStorage, cfg, msg.Value)
			continue
		}

		log.Printf("Error unmarshalling message: %v", err)
	}
}

func processOrders(orders []map[string]interface{}, cache *cache.Cache, dbStorage *storage.Storage, cfg *config.KafkaConfig, data []byte) {
	for _, order := range orders {
		processOrder(order, cache, dbStorage, cfg, data) // Передаем оригинальные данные
	}
}

func processOrder(order map[string]interface{}, cache *cache.Cache, dbStorage *storage.Storage, cfg *config.KafkaConfig, data []byte) {
	orderUID, ok := order["order_uid"].(string)
	if !ok || orderUID == "" {
		log.Println("Invalid order format: missing order_uid")
		return
	}

	cache.Set(orderUID, data)
	log.Printf("Caching order %s", orderUID)

	if err := saveWithRetry(context.Background(), dbStorage, orderUID, data, int(cfg.Retries), cfg.RetryDelay); err != nil {
		log.Printf("Error saving order to DB after retries: %v", err)
	}
}
