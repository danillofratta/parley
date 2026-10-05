package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr      string
	DatabaseURL   string
	WebhookSecret string
	TenantID      string

	KafkaBrokers         []string
	OutboxPollInterval   time.Duration
	OutboxBatchSize      int
	OutboxPublishTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr: getenv("HTTP_ADDR", ":8081"),
		// The service connects with its own user, which can only touch the gateway schema (ADR 0006).
		DatabaseURL:   getenv("GATEWAY_DATABASE_URL", "postgres://gateway:gateway@localhost:5432/parley?sslmode=disable"),
		WebhookSecret: os.Getenv("TELEGRAM_WEBHOOK_SECRET"),
		TenantID:      getenv("DEFAULT_TENANT_ID", "00000000-0000-0000-0000-000000000001"),
		KafkaBrokers:  strings.Split(getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"), ","),
	}
	if cfg.WebhookSecret == "" {
		return Config{}, errors.New("TELEGRAM_WEBHOOK_SECRET is required")
	}

	var err error
	if cfg.OutboxPollInterval, err = durationEnv("OUTBOX_POLL_INTERVAL", 500*time.Millisecond); err != nil {
		return Config{}, err
	}
	if cfg.OutboxPublishTimeout, err = durationEnv("OUTBOX_PUBLISH_TIMEOUT", 10*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.OutboxBatchSize, err = intEnv("OUTBOX_BATCH_SIZE", 100); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration such as 500ms", key)
	}
	return d, nil
}

func intEnv(key string, fallback int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return n, nil
}
