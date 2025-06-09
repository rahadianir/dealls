package logger

import (
	"log/slog"
	"os"
)

func InitLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.String("app", "hrapp"),
	)
}
