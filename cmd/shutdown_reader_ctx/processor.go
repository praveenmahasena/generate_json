// Package main ...
package main

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

func processAndDelete(state directoryReaderInterface, fileRead, bytesRead *atomic.Uint64) error {
	// I want to do something like this
	// state.(directoryReader) and do the other process
	// what if you make this func to recieve any other struct that satisfy directoryReaderInterface other than directoryReader in the future?
	// I'll use pkg.go.dev/reflect i know it's going to be doing type caste on runtime and this might slow up idk maybe I'm wrong please back me on that
	s := reflect.ValueOf(state)
	fmt.Println(s)

	return nil
}
