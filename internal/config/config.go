package config

import (
	"context"
	"log"
	"os"
	"strconv"
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
			Name:    getEnvString("APP_NAME", "hr-app"),
			Version: getEnvString("APP_VERSION", "1.0.0"),
			IsDebug: getEnvBool("IS_DEBUG_MODE", false),

			ExpiryTime:   getEnvDuration("ACCESS_TOKEN_EXPIRY_DURATION", "10h"),
			JWTSecretKey: getEnvString("JWT_SECRET_KEY", "secret"),
		},
		DB: &DB{
			URL:     getEnvString("DB_URL", "postgres://postgres:postgres@localhost:5432/database?sslmode=disable"),
			MaxConn: getEnvInt("DB_MAX_CONN", 20),
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

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	finalVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return finalVal
}

func getEnvDuration(key string, defaultVal string) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		dur, err := time.ParseDuration(defaultVal)
		if err != nil {
			log.Fatal("failed to parse ", key, "duration config: ", err)
		}

		return dur
	}

	dur, err := time.ParseDuration(val)
	if err != nil {
		log.Fatal("failed to parse ", key, "duration config: ", err)
	}

	return dur
}
