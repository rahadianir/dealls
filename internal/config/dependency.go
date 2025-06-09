package config

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type CommonDependencies struct {
	Config *Config
	DB     *sqlx.DB
	Logger *slog.Logger
}
