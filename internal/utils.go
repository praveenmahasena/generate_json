package internal

import (
	"log/slog"
	"os"
)

type Js struct {
	ID uint `json:"id"`
}

func NewLogger(w *os.File, source bool, level int) *slog.Logger {
	options := slog.HandlerOptions{
		AddSource: source,
		Level:     slog.Level(level),
	}
	handler := slog.NewJSONHandler(os.Stdout, &options)
	return slog.New(handler)
}
