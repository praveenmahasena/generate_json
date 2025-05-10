package internal

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/praveenmahasena/generate_json/internal/cli"
	"github.com/praveenmahasena/generate_json/internal/jsonreader"
	"github.com/praveenmahasena/generate_json/internal/jsonwriter"
)

// pretty sure logger is no more needed since I took off go routine
// 15 mins later --
// actually I'm wrong this is needed
func logger() *slog.Logger {
	w := io.MultiWriter(os.Stdout)
	slogOpt := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.Level(1),
	}
	l := slog.NewJSONHandler(w, &slogOpt)
	return slog.New(l)
}

func Run() error {
	amount := cli.ReadAmount()
	js := jsonwriter.New(amount, logger())
	return js.Generate()
}

func Read() error {
	fileRead, byteRead, err := jsonreader.Read(logger())
	if err != nil {
		return err
	}
	fmt.Printf("total files read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", byteRead)
	return nil
}
