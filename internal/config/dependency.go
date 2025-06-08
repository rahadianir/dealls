package config

import "github.com/jackc/pgx/v5/pgxpool"

type CommonDependencies struct {
	DBPool *pgxpool.Pool
}
