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

		var order map[string]interface{}
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error unmarshalling order: %v", err)
			continue
		}

		orderUID, ok := order["order_uid"].(string)
		if !ok || orderUID == "" {
			log.Println("Invalid order format: missing order_uid")
			continue
		}

		cache.Set(orderUID, msg.Value)

		if err := saveWithRetry(context.Background(), dbStorage, orderUID, msg.Value, int(cfg.Retries), cfg.RetryDelay); err != nil {
			log.Printf("Error saving order to DB after retries: %v", err)
		}
	}
}
