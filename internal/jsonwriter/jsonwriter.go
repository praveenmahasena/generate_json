package jsonwriter

import (
	"encoding/json"
	"log/slog"
	"fmt"
	"path"
	"sync"
	"syscall"
)

type Jsonwriter struct {
	Amount uint
	l      *slog.Logger
}

type Js struct {
	ID uint `json:"id"`
}

func New(id uint,l *slog.Logger) *Jsonwriter {
	return &Jsonwriter{id,l}
}

func (j *Jsonwriter) Generate() error {
	dir, dirErr := syscall.Getwd()
	if dirErr != nil {
		return fmt.Errorf("error during getting work dir %+v", dirErr)
	}
	path := path.Join(dir, "/json/")                                                // I could've done dir+"./json/" but js using this for fancy
	if err := syscall.Unlink(path); !(err.Error() == "no such file or directory") { // at this point I should implement a errors package for myself for these kinda case
		return fmt.Errorf("error during removing existing pre generated json dir %+v ", err)
	}
	if err := syscall.Mkdir(path, 0700); err != nil {
		return fmt.Errorf("error during creating json dir %+v", err)
	}
	if err := syscall.Chdir(path); err != nil {
		return fmt.Errorf("error during changing to json dir %+v", err)
	}
	wg := &sync.WaitGroup{}
	for i := range j.Amount {
		wg.Add(1)
		go writeFile(i, wg,j.l)
	}
	wg.Wait()
	return nil
}

func writeFile(i uint, wg *sync.WaitGroup,l *slog.Logger) {
	// why logger for error handling?
	// well these are go routines and they should not return anything
	// I can do a channel and handle errors but at this point I don't see a value
	// also in my opinion it would cost a lot since we are generating massive amounts of files and handling json
	// so I prefer to have them logged to stdout
	defer wg.Done()
	fName := fmt.Sprintf("%v.json", i)
	// well this is a bug
	// i want to create a read, deletable file that's all
	// but this one creates a executable
	// which makes my *delete previously generated feature* fail
	// someone help AHHHHH
	// or maybe the bug is in line : 33
	file, fileErr := syscall.Creat(fName, 0600)

	if fileErr != nil {
		e := fmt.Sprintf("error during creating %v", fName)
		// I do not like to do this error unwraping
		// but for now it's alright
		l.Error(e,fileErr.Error(),"")
		return
	}
	defer syscall.Close(file)
	content := Js{i}
	b, bErr := json.MarshalIndent(content, "", "")

	if bErr != nil {
		e := fmt.Sprintf("error during generating json for file %v", fName)
		l.Error(e,bErr.Error(),"")
		return
	}
	_, err := syscall.Write(file, b)
	if err != nil {
		e := fmt.Sprintf("error during writing in generating json into file %v", fName)
		l.Error(e,err.Error(),"")
		return
	}
}
