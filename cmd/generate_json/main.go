package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/praveenmahasena/generate_json/internal"
)

var (
	defaultAmount      uint = 100
	defaultDirectories uint = 100
)

type Jsonwriter struct {
	amount      uint
	directories uint
}

func main() {
	logger := internal.NewLogger(os.Stderr, true, 1)
	amount, directories := readAmount()
	js := new(amount, directories)
	if err := js.generate(logger); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func new(id, directories uint) *Jsonwriter {
	return &Jsonwriter{id, directories}
}

func readAmount() (uint, uint) {
	amount := flag.Uint("amount", defaultAmount, "amount of json files this binary should generate would be defaulted to 100")
	directories := flag.Uint("directories", defaultDirectories, "amount of json directories this binary should generate would be defaulted to 100")
	flag.Parse()
	return *(amount), *(directories) // call me old idc but been doing Clang for awhile and this *() reference seems more appealing then just doing *
}

func (j *Jsonwriter) generate(l *slog.Logger) error {
	if err := os.RemoveAll("json"); err != nil {
		return fmt.Errorf("error during removing previous json dir %+v ", err)
	}
	if err := os.Mkdir("json", 0777); err != nil {
		return fmt.Errorf("error during making dir for json %+v", err)
	}
	if err := os.Chdir("./json"); err != nil {
		return fmt.Errorf("error during moving into dir for json with value %v ,%+v", "json", err)
	}
	for d := range j.directories {
		dirName := fmt.Sprintf("%vjson", d)
		if err := os.Mkdir(dirName, 0777); err != nil {
			l.Error("error during making dir for json", "error value", err.Error(), "process", "skipping dir", "progress", fmt.Sprintf("%v", dirName))
			continue
		}
		if err := os.Chdir("./" + dirName); err != nil {
			l.Error("error during changing dir for json", "error value", err.Error(), "process", "skipping dir", "progress", fmt.Sprintf("%v", dirName))
			continue
		}

		for a := range j.amount {
			if err := writeFile(a); err != nil {
				l.Error("error during writing into json file", "error value", err.Error(), "process", "retry 1", "progress", fmt.Sprintf("writing again into %v", dirName))
				if err := writeFile(a); err != nil {
					l.Error("error during writing into json file", "error value", err.Error(), "process", "retry 1 failed", "progress", fmt.Sprintf("skipping file %v", dirName))
				}
			}
		}
		if err := os.Chdir("../"); err != nil {
			return fmt.Errorf("error during moving moving dir in json %+v", err)
		}
	}
	return nil
}

func writeFile(i uint) error {
	fName := fmt.Sprintf("%v.json", i)
	js := internal.Js{ID: i}
	b, err := json.MarshalIndent(js, "", "")
	if err != nil {
		return fmt.Errorf("error during preparing json content for file %v with value %+v", fName, err)
	}
	if err := os.WriteFile(fName, []byte(b), 0666); err != nil {
		return fmt.Errorf("error during writing into file %v with value %+v", fName, err)
	}
	return nil
}
