package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rahadianir/dealls/internal/attendance"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/middleware"
	"github.com/rahadianir/dealls/internal/payroll"
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

	err = db.Ping()
	if err != nil {
		logger.ErrorContext(ctx, "failed to ping db connection", slog.Any("error", err))
		os.Exit(1)
	}

	deps := config.CommonDependencies{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}

	// init http routes
	routes := initRoutes(&deps)

	// setup server
	var srv http.Server

	// setup for graceful shutdown
	idleConnectionClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shutting down.
		logger.InfoContext(ctx, "HTTP Server is shutting down")
		db.Close()
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			logger.ErrorContext(ctx, "HTTP Server fails to shutting down", slog.Any("error", err))
		}
		close(idleConnectionClosed)

	}()

	srv.Addr = fmt.Sprintf(":%d", 8080)
	srv.Handler = routes

	logger.InfoContext(ctx, "HTTP Server starts!")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.ErrorContext(ctx, "HTTP Server fails to start listen and serve", slog.Any("error", err))
	}

	<-idleConnectionClosed
	logger.InfoContext(ctx, "Bye!")
}

func initRoutes(deps *config.CommonDependencies) http.Handler {
	// wiring layers
	// shared packages
	jwtHelper := xjwt.XJWT{}

	// repository
	userRepo := user.NewUserRepository(deps)
	attRepo := attendance.NewAttendanceRepository(deps)
	payrollRepo := payroll.NewPayrollRepository(deps)

	// logic
	userLogic := user.NewUserLogic(deps, *userRepo, &jwtHelper)
	attLogic := attendance.NewAttendanceLogic(deps, *attRepo)
	payrollLogic := payroll.NewPayrollLogic(deps, *payrollRepo, *userLogic, *attRepo)

	// handler
	userHandler := user.NewUserHandler(deps, *userLogic)
	attHandler := attendance.NewAttendanceHandler(deps, *attLogic)
	payrollHandler := payroll.NewPayrollHandler(deps, *payrollLogic)

	// setup middlewares
	authMW := middleware.NewAuthMiddleware(deps, &jwtHelper)
	traceMW := middleware.TracerMiddleware{}
	r := chi.NewRouter()

	r.Use(traceMW.Tracer)

	r.Post("/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(authMW.AuthOnly) // check whether the user is logged in with proper auth and embed user id in context
		r.Post("/attendance", attHandler.SubmitAttendance)
		r.Post("/overtime", attHandler.SubmitOvertime)
		r.Post("/reimbursement", attHandler.SubmitReimbursement)
		r.Post("/payroll/period", payrollHandler.SetPayrollPeriod)
	})

	return r
}
