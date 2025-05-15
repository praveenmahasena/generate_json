package logger

import (
	"io"
	"log/slog"
	"os"
)

func Logger() *slog.Logger {
	w := io.MultiWriter(os.Stdout)
	slogOpt := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.Level(1),
	}
	l := slog.NewJSONHandler(w, &slogOpt)
	return slog.New(l)
}
