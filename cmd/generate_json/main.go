package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/praveenmahasena/generate_json/internal"
)

var (
	defaultAmount uint = 10_000
)

type Jsonwriter struct {
	Amount uint
}

func main() {
	amount := readAmount()
	js := new(amount)
	if err := js.generate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func new(id uint) *Jsonwriter {
	return &Jsonwriter{id}
}

func readAmount() uint {
	amount := flag.Uint("amount", defaultAmount, "amount of json files this binary should generate would be defaulted to 10_000")
	flag.Parse()
	return *(amount) // call me old idc but been doing Clang for awhile and this *() reference seems more appealing then just doing *
}

func (j *Jsonwriter) generate() error {
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		return fmt.Errorf("error during getting work dir %+v", dirErr)
	}
	path := dir + "/json/"
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("error during removing existing pre generated json dir %+v ", err)
	}
	if err := os.Mkdir(path, 0777); err != nil {
		return fmt.Errorf("error during creating json dir %+v", err)
	}
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("error during changing to json dir %+v", err)
	}
	for i := range j.Amount {
		if err := writeFile(i); err != nil {
			log.Printf("file write error %v", err)
		}
	}
	return nil
}

func writeFile(i uint) error {
	fileName := fmt.Sprintf("%v.json", i)
	content, err := json.MarshalIndent(internal.Js{i}, " ", "")
	if err != nil {
		return fmt.Errorf("error during generating json %+v", err)
	}
	fileErr := os.WriteFile(fileName, content, 0644)
	if fileErr != nil {
		return fmt.Errorf("error during writing into json file %+v", fileErr)
	}
	return nil
}
