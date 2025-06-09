package config

import (
	"context"
	"os"
	"time"
)

type Config struct {
	App *App
	DB  *DB
}

type App struct {
	// general app related config
	Name    string
	Version string

	// auth
	ExpiryTime time.Duration
	SecretKey  string
}
type DB struct {
	URL     string
	MaxConn int
	// add more db connection config
}

func InitConfig(ctx context.Context) *Config {

	return &Config{
		App: &App{
			Name:       "hr-app",
			Version:    "1.0.0",
			ExpiryTime: time.Duration(5 * time.Minute),
			SecretKey:  "secret",
		},
		DB: &DB{
			URL:     getEnvString("", "postgres://postgres:postgres@localhost:5432/database?sslmode=disable"),
			MaxConn: 20,
		},
	}
}

func getEnvString(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}
