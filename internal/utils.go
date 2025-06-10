package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Js struct {
	ID uint `json:"id"`
}

type JobReader interface {
	Read(ctx context.Context, c chan<- string)
	ProcessAndDelete(fileName string) error
}

func NewLogger(w *os.File, source bool, level int) *slog.Logger {
	options := slog.HandlerOptions{
		AddSource: source,
		Level:     slog.Level(level),
	}
	handler := slog.NewJSONHandler(w, &options)
	return slog.New(handler)
}

// ShutDown a general ShutDown func
func ShutDown(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
}
