package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rahadianir/dealls/internal/config"
)

func StartServer() {
	ctx := context.Background()

	// setup config
	cfg := config.InitConfig(ctx)

	// init common dependencies
	// init database connection pool
	pgConfig, err := pgxpool.ParseConfig(cfg.DB.URL)
	if err != nil {

	}
	pgConfig.MaxConns = int32(cfg.DB.MaxConn)

	dbPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {

	}
	defer dbPool.Close()

	deps := config.CommonDependencies{
		DBPool: dbPool,
	}

	// serve http
	serveHTTP(ctx, &deps)
}

func serveHTTP(ctx context.Context, deps *config.CommonDependencies) {}
