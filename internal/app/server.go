package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rahadianir/dealls/internal/attendance"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/middleware"
	"github.com/rahadianir/dealls/internal/pkg/logger"
	"github.com/rahadianir/dealls/internal/pkg/xjwt"
	"github.com/rahadianir/dealls/internal/user"
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
		Config: cfg,
		DB:     db,
		Logger: logger,
	}

	// serve http
	routes := initRoutes(ctx, &deps)

	logger.InfoContext(ctx, "server starts")
	err = http.ListenAndServe(":8080", routes)
	if err != nil {
		panic(err)
	}
}

func initRoutes(ctx context.Context, deps *config.CommonDependencies) http.Handler {
	// wiring layers
	// shared packages
	jwtHelper := xjwt.XJWT{}

	// repository
	userRepo := user.NewUserRepository(deps)
	attRepo := attendance.NewAttendanceRepository(deps)

	// logic
	userLogic := user.NewUserLogic(deps, *userRepo, &jwtHelper)
	attLogic := attendance.NewAttendanceLogic(deps, *attRepo)

	// handler
	userHandler := user.NewUserHandler(deps, *userLogic)
	attHandler := attendance.NewAttendanceHandler(deps, *attLogic)

	r := chi.NewRouter()

	traceMW := middleware.TracerMiddleware{}
	r.Use(traceMW.Tracer)

	r.Post("/login", userHandler.Login)

	r.Post("/attendance", attHandler.SubmitAttendance)

	return r
}
