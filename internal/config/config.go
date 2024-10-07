package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Kafka    KafkaConfig    `yaml:"kafka"`
	Cache    CacheConfig    `yaml:"cache"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type ServerConfig struct {
	Env     string        `yaml:"env" env-default:"development"`
	Port    string        `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type DatabaseConfig struct {
	DatabaseURL string `yaml:"database_url" env-required:"true"`
}

type KafkaConfig struct {
	Broker     string        `yaml:"broker" env-required:"true"`
	Topic      string        `yaml:"topic" env-required:"true"`
	GroupID    string        `yaml:"group_id" env-required:"true"`
	Retries    uint          `yaml:"retries" env-default:"3"`
	RetryDelay time.Duration `yaml:"retry_delay" env-default:"1s"`
}

type CacheConfig struct {
	Expiration      time.Duration `yaml:"expiration" env-default:"10m"`
	CleanupInterval time.Duration `yaml:"cleanup_interval" env-default:"5m"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env-default:"info"`
}

func MustLoad(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
		return nil, err
	}

	return &cfg, nil
}
