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
	errCh := make(chan error)
	// why context here I could've done a (chan bool) and make this whole thing work in a way
	// but using it cuz it's more of a design pattern or convension
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()     // no context leak
	defer close(errCh) // no channel leak

	go jsonreader.GracefulRead(ctx, errCh,logger())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT,syscall.SIGTERM)

	select {
	case <-sigCh:
		fmt.Println("cancelling...")
		cancel()
	case err := <-errCh:
		return err
	}

	// I'm pretty sure after the cancel() no files would open
	// which means if no file open no file gets json unmarshed
	// there are no way i'm gonna have any kind of error
	// but also I gotta make this sync or it would just suddenly stop
	// no graceful shutdown concept being followed
	// I can do sync.WaitGroup{}
	// but in this case I think it's more easy to handle it with errCh since it has been already created for err handling
	// and pretty sure there this won't give out any err after <cmd+c> pressed

	<-errCh
	return nil
}
