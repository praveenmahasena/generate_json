package jsonwriter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"
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
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		return fmt.Errorf("error during getting work dir %+v", dirErr)
	}
	path := path.Join(dir, "/json/") // I could've done dir+"./json/" but js using this for fancy
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
			j.l.Error("file write error", err.Error(), "")
		}
	}
	return nil
}

func writeFile(i uint) error {
	fileName := fmt.Sprintf("%v.json", i)
	content, err := json.MarshalIndent(Js{i}, " ", "")
	if err != nil {
		return fmt.Errorf("error during generating json %+v", err)
	}
	fileErr := os.WriteFile(fileName, content, 0644)
	if fileErr != nil {
		return fmt.Errorf("error during writing into json file %+v", fileErr)
	}
	return nil
}
