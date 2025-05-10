package internal

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

// actually I'll be using go routine here
// but not to read files
// my brain is too tiny I cannot think of other ways :)
// also channels for synced communication or data passing
func GracefulRead() error {
	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // let's face it // I don't wanna have embarasing context leaks no fun for me

	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTRAP)

	go func(ctx context.Context,err chan error) {
		jsonreader.GracefulRead(ctx,logger(),err)
	}(ctx,errCh)

	<-sigCh
	cancel()

	return <-errCh
}
