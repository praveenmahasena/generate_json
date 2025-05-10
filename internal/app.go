package internal

import (
	"io"
	"log/slog"
	"os"

	"github.com/praveenmahasena/generate_json/internal/cli"
	"github.com/praveenmahasena/generate_json/internal/jsonwriter"
)

func Run() error {
	amount := cli.ReadAmount()
	w := io.MultiWriter(os.Stdout)
	slogOpt := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.Level(1),
	}
	l := slog.NewJSONHandler(w, &slogOpt)
	sloger := slog.New(l)
	js := jsonwriter.New(amount, sloger)
	return js.Generate()
}
