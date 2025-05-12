package internal

import (
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
	sigCh := make(chan os.Signal,1)
	signal.Notify(sigCh,syscall.SIGINT,syscall.SIGTERM)
	fileRead, byteRead, err := jsonreader.GracefulRead(sigCh,logger())
	if err!=nil{
		return nil
	}
	fmt.Printf("total files read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", byteRead)
	return nil
}
