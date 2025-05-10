package jsonwriter

import (
	"fmt"
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"syscall"
)

type Jsonwriter struct {
	Amount uint
	l      *slog.Logger
}

type Js struct {
	ID uint `json:"id"`
}

func New(id uint, l *slog.Logger) *Jsonwriter {
	return &Jsonwriter{id, l}
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
	for i := range j.Amount {
		if err := writeFile(i); err != nil {
			j.l.Error("file write error", err.Error(), "")
		}
	}
	return nil
}

func writeFile(i uint) error {
	fileName:=fmt.Sprintf("%v.json",i)
	content,err:=json.Marshal(Js{i})
	if err!=nil{
		return fmt.Errorf("error during generating json %+v",err)
	}
	fileErr:=os.WriteFile(fileName,content,os.FileMode(os.O_EXCL))
	if fileErr!=nil{
		return fmt.Errorf("error during writing into json file %+v",fileErr)
	}
	return nil
}
