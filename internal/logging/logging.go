package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type LoggerConfig struct {
	Level      string
	Format     string
	WithSource bool
}

type ctxKey struct{}

func New(cfg LoggerConfig) *slog.Logger {
	var lv slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		lv = slog.LevelDebug
	case "info":
		lv = slog.LevelInfo
	case "warn", "warning":
		lv = slog.LevelWarn
	case "error":
		lv = slog.LevelError
	default:
		lv = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level:     lv,
		AddSource: cfg.WithSource,
	}
	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "plain":
		handler = NewPlainHandler(os.Stdout, opts.Level).
			WithSource(cfg.WithSource).
			WithTimestamp(lv < slog.LevelInfo)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}

func ToContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(ctxKey{}).(*slog.Logger)
	if !ok || l == nil {
		return slog.Default()
	}
	return l
}
