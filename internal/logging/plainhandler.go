package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"time"
)

type PlainHandler struct {
	w          io.Writer
	level      slog.Leveler
	withSource bool
}

func NewPlainHandler(w io.Writer, level slog.Leveler, withSource bool) *PlainHandler {
	return &PlainHandler{w: w, level: level, withSource: withSource}
}

func (h *PlainHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level.Level()
}

func (h *PlainHandler) Handle(_ context.Context, r slog.Record) error {
	ts := r.Time.Format(time.RFC3339)
	level := r.Level.String()

	// base message
	line := fmt.Sprintf("%s %-5s %s", ts, level, r.Message)

	// optional source
	if h.withSource && r.PC != 0 {
		fs := r.Source()
		line = fmt.Sprintf("%s (%s:%d)", line, filepath.Base(fs.File), fs.Line)
	}

	_, err := fmt.Fprintln(h.w, line)
	return err
}

func (h *PlainHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *PlainHandler) WithGroup(name string) slog.Handler       { return h }
