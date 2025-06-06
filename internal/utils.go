package internal

import (
	"log/slog"
	"os"
	"context"
)

type Js struct {
	ID uint `json:"id"`
}

type JobReader interface {
  Read(ctx context.Context, c chan<- string)
}

func NewLogger(w *os.File, source bool, level int) *slog.Logger {
	options := slog.HandlerOptions{
		AddSource: source,
		Level:     slog.Level(level),
	}
	handler := slog.NewJSONHandler(os.Stdout, &options)
	return slog.New(handler)
}
