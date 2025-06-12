package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// Js structs to hold json structure
type Js struct {
	ID uint `json:"id"`
}

// NewLogger helper method for logging out
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

// GetPath ...
func GetPath() (string, error) {
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return "", fmt.Errorf("error during getting working dir with value %#v", wdErr)
	}
	return wd + "/json", nil
}
