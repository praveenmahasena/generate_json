package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/praveenmahasena/generate_json/internal/jsonreader"
	"github.com/praveenmahasena/generate_json/internal/logger"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	fileRead, byteRead, err := jsonreader.GracefulRead(sigCh,logger.Logger())

	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		// I did not wanted to return btw
		// just print the error and then
		// show out that it only read 0 values
	}
	fmt.Printf("total files read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", byteRead)
}

