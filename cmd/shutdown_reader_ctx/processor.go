// Package main ...
package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/praveenmahasena/generate_json/internal"
)

func processAndDelete(state directoryReaderInterface, nameCh <-chan string, wg *sync.WaitGroup) {
	// I want to do something like this
	// state.(directoryReader) and do the other process
	// what if you make this func to recieve any other struct that satisfy directoryReaderInterface other than directoryReader in the future?
	// I'll use pkg.go.dev/reflect i know it's going to be doing type caste on runtime and this might slow up idk maybe I'm wrong please back me on that

	defer wg.Done()

	for name := range nameCh {
		b, err := os.ReadFile(name)
		if err != nil {
			log.Printf("error during reading out file %v with value %#v", name, err)
			continue
		}
		if err := json.Unmarshal(b, &internal.Js{}); err != nil {
			log.Printf("error during reading out to json file %v with value %#v", name, err)
			continue
		}
		state.addState(uint64(len(b)), 1)
		if err := os.Remove(name); err != nil {
			log.Printf("error during removing file %v with value of % #v", name, err)
		}
	}
}
