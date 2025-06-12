package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"

	"github.com/praveenmahasena/generate_json/internal"
)

func main() {
	path, pathErr := internal.GetPath()
	if pathErr != nil {
		fmt.Fprintln(os.Stderr, pathErr)
		os.Exit(-1)
	}
	d := internal.NewdirectoryReader(path, &atomic.Uint64{}, &atomic.Uint64{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go internal.ShutDown(cancel)
	nameCh := make(chan string)
	wg := &sync.WaitGroup{}
	go processAndDelete(d, nameCh, wg)

	if err := d.Read(ctx, nameCh); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	wg.Wait()
	d.ShowStats()
}

func processAndDelete(state internal.DirectoryReaderInterface, nameCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for name := range nameCh {
		d, err := os.ReadFile(name)
		if err != nil {
			log.Printf("file read error on %v with value %#v", name, err)
			continue
		}
		if err := json.Unmarshal(d, &internal.Js{}); err != nil {
			log.Printf("error during parsing file %v into json struct with value %#v", name, err)
			continue
		}
		state.AddState(uint64(len(d)), 1)
	}
}
