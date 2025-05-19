package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	defaultAmount      uint = 100
	defaultDirectories uint = 1000
)

type Jsonwriter struct {
	amount      uint
	directories uint
}

func main() {
	amount, directories := readAmount()
	js := new(amount, directories)
	if err := js.generate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func new(id, directories uint) *Jsonwriter {
	return &Jsonwriter{id, directories}
}

func readAmount() (uint, uint) {
	amount := flag.Uint("amount", defaultAmount, "amount of json files this binary should generate would be defaulted to 10_000")
	directories := flag.Uint("directories", defaultDirectories, "amount of json directories this binary should generate would be defaulted to 1")
	flag.Parse()
	return *(amount), *(directories) // call me old idc but been doing Clang for awhile and this *() reference seems more appealing then just doing *
}

func (j *Jsonwriter) generate() error {
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
			log.Printf("error during making dir for json %+v", err)
			continue
		}
		if err := os.Chdir("./" + dirName); err != nil {
			log.Printf("error during moving into dir for json with value %v ,%+v", dirName, err)
			continue
		}

		for a := range j.amount {
			if err := writeFile(a); err != nil {
				log.Println()
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
	b := `{
		"id":`+fmt.Sprintf("%v",i)+`
	}`
	if err := os.WriteFile(fName, []byte(b), 0666); err != nil {
		return fmt.Errorf("error during writing into file %v with value %+v", fName, err)
	}
	return nil
}
