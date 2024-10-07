package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"L0/internal/config"
	"github.com/segmentio/kafka-go"
)

const (
	configPath = "config/config.yaml"
	modelPath  = "data/model.json"
)

func main() {
	cfg, err := config.MustLoad(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	data, err := os.ReadFile(modelPath)
	if err != nil {
		log.Fatalf("Error reading model.json: %v", err)
	}

	writer := kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Broker),
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("order_uid"),
			Value: data,
		},
	)
	if err != nil {
		log.Fatalf("Error writing message to Kafka: %v", err)
	}

	fmt.Println("Order sent successfully")

	writer.Close()
}
