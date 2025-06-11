// Package main ...
package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/praveenmahasena/generate_json/internal"
)

func main() {
	path, pathErr := getPath()
	if pathErr != nil {
		fmt.Fprintln(os.Stdout, pathErr)
		syscall.Exit(-1)
	}
	// I do not like this way of handling but since pkg.go.dev.sync/atomic go routine safe I'm taking this approach
	fileRead, bytesRead := &atomic.Uint64{}, &atomic.Uint64{}
	s := newState(path, fileRead, bytesRead)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go internal.ShutDown(cancel)
	nameCh := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go processAndDelete(s, nameCh, wg)
	if err := s.read(ctx, nameCh); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	close(nameCh)
	wg.Wait()
	s.showStats()
}
