package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/logger"
)

func StartServer() {
	ctx := context.Background()

	// setup config
	cfg := config.InitConfig(ctx)

	// init common dependencies
	// init logger
	logger := logger.InitLogger()

	// init database connection pool
	db, err := sqlx.Open("postgres", cfg.DB.URL)
	if err != nil {
		logger.ErrorContext(ctx, "failed to open db connection", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	deps := config.CommonDependencies{
		DB:     db,
		Logger: logger,
	}

	// serve http
	serveHTTP(ctx, &deps)
}

func serveHTTP(ctx context.Context, deps *config.CommonDependencies) {}
