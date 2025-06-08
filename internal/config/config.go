package config

import "context"

type Config struct {
	DB *DB
}

type DB struct {
	URL     string
	MaxConn int
}

func InitConfig(ctx context.Context) *Config {

	return &Config{}
}
