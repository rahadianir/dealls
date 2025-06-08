package config

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type CommonDependencies struct {
	DB *sqlx.DB
	Logger *slog.Logger
}
