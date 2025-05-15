package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/praveenmahasena/generate_json/internal/jsonreader"
	"github.com/praveenmahasena/generate_json/internal/logger"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	fileRead, byteRead, err := jsonreader.GracefulRead(sigCh, logger.Logger())
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	fmt.Println(fileRead, byteRead)
}
