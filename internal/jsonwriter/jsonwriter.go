package jsonwriter

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"sync"
	"syscall"
)

type Jsonwriter struct {
	Amount uint
}

type Js struct {
	ID uint `json:"id"`
}

func New(id uint) *Jsonwriter {
	return &Jsonwriter{id}
}

func (j *Jsonwriter) Generate() error {
	dir, dirErr := syscall.Getwd()
	if dirErr != nil {
		return fmt.Errorf("error during getting work dir %+v", dirErr)
	}
	path := path.Join(dir, "/json/")                                               // I could've done dir+"./json/" but js using this for fancy
	if err := syscall.Unlink(path); !(err.Error() == "no such file or directory") { // at this point I should implement a errors package for myself for these kinda case
		return fmt.Errorf("error during removing existing pre generated json dir %+v :", err)
	}
	if err := syscall.Mkdir(path, 0777); err != nil {
		return fmt.Errorf("error during creating json dir %+v", err)
	}
	if err := syscall.Chdir(path); err != nil {
		return fmt.Errorf("error during changing to json dir %+v", err)
	}
	wg := &sync.WaitGroup{}
	for i := range j.Amount {
		wg.Add(1)
		go writeFile(i, wg)
	}
	wg.Wait()
	return nil
}

func writeFile(i uint, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("error here?")
	fName := fmt.Sprintf("%v.json", i)
	file, fileErr := syscall.Creat(fName, syscall.O_WRONLY)

	if fileErr != nil {
		e := fmt.Errorf("error during creating %v with %+v", fName, fileErr)
		log.Println(e)
		return
	}
	defer syscall.Close(file)
	content := Js{i}
	b, bErr := json.MarshalIndent(content, "", "")

	if bErr != nil {
		e := fmt.Errorf("error during generating json for  %v with %+v", fName, bErr)
		log.Println(e)
		return
	}
	_, err := syscall.Write(file, b)
	if err != nil {
		e := fmt.Errorf("error during writing into json file %v with %+v", fName, err)
		log.Println(e)
		return
	}
}
