package config

import (
	"context"
	"os"
	"strings"
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
	IsDebug bool

	// auth
	ExpiryTime   time.Duration
	JWTSecretKey string
}
type DB struct {
	URL     string
	MaxConn int
	// add more db connection config
}

func InitConfig(ctx context.Context) *Config {

	return &Config{
		App: &App{
			Name:         "hr-app",
			Version:      "1.0.0",
			IsDebug:      getEnvBool("IS_DEBUG_MODE", false),

			ExpiryTime:   time.Duration(10 * time.Minute),
			JWTSecretKey: "secret",
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

func getEnvBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)

	switch strings.ToLower(val) {
	case "1", "true", "on":
		return true
	default:
		return defaultValue
	}
}
