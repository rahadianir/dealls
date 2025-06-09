package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/rahadianir/dealls/internal/config"
)

type CustomHandler struct {
	handler slog.Handler
}

func InitLogger() *slog.Logger {
	return slog.New(&CustomHandler{
		handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.MessageKey {
					a.Key = "message"
				}

				if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
				}

				return a
			},
		}),
	})
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := ctx.Value(config.RequestIDKey).(string); ok {
		r.AddAttrs(slog.String("request_id", reqID))
	}
	return h.handler.Handle(ctx, r)
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithGroup(name),
	}
}
