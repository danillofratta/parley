package config

import (
	"os"
	"errors"
)

type Config struct {
	HTTPAddr string
	DatabaseURL string
	WebhookSecret string
	TenantId string
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:      getenv("HTTP_ADDR", ":8081"),
		DatabaseURL:   getenv("GATEWAY_DATABASE_URL", "postgres://gateway:gateway@localhost:5432/parley"),
		WebhookSecret: os.Getenv("TELEGRAM_WEBHOOK_SECRET"),
		TenantID:      getenv("DEFAULT_TENANT_ID", "00000000-0000-0000-0000-000000000001"),
	}
	if cfg.WebhookSecret == "" {
		return Config{}, errors.New("TELEGRAM_WEBHOOK_SECRET is required")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}