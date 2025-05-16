package jsonreader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/praveenmahasena/generate_json/internal/jsondelete"
	"github.com/praveenmahasena/generate_json/internal/jsonwriter"
)

type (
	filesRead uint
	bytesRead uint
)

func Read(l *slog.Logger) (filesRead, bytesRead, error) {
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return 0, 0, fmt.Errorf("error during getting work dir :%+v", wdErr)
	}
	if err := os.Chdir(wd + "/json/"); err != nil {
		return 0, 0, fmt.Errorf("error during getting into json work dir :%+v", err)
	}
	files, filesErr := os.Open(".")
	if filesErr != nil {
		return 0, 0, fmt.Errorf("error during opening json work dir :%+v", filesErr)
	}
	defer files.Close()
	fileNames, fileNameErr := files.Readdirnames(10_000_000_000) // and no this one is not gonna create an array with cap string 10M
	// let's break this
	// the fileNames is gonna have a arrayslice is gonna hold all the file names in the dir
	// not being sorted
	if fileNameErr != nil {
		return 0, 0, fmt.Errorf("error during opening json work dir :%+v", filesErr)
	}
	var (
		fileRead  filesRead = 0
		byteRead  bytesRead = 0
		fileBytes           = make([]byte, 50)
		file      *os.File
		err       error
		n         int
	)
	for _, name := range fileNames {
		file, err = os.Open(name)
		if err != nil {
			log.Printf("error during opening file %v with value %+v", name, err)
			file.Close()
			continue
		}
		n, err = io.ReadFull(file, fileBytes)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			log.Printf("error during reading file %v with value %+v", name, err)
			file.Close()
			continue
		}
		fileRead += 1
		byteRead += bytesRead(n)
		file.Close()
		js := &jsonwriter.Js{}
		if err = json.Unmarshal(fileBytes[:n], js); err != nil {
			log.Printf("error during unmarshal file %v with value %+v", name, err)
		}
		jsondelete.DeleteFile(name)
	}
	return fileRead, byteRead, nil
}

// fileRead, byteRead, error
func GracefulRead(sigCh chan os.Signal, l *slog.Logger) (filesRead, bytesRead, error) {
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return 0, 0, fmt.Errorf("error during getting work dir :%+v", wdErr)
	}
	if err := os.Chdir(wd + "/json/"); err != nil {
		return 0, 0, fmt.Errorf("error during getting into json work dir :%+v", err)
	}
	files, filesErr := os.Open(".")
	if filesErr != nil {
		return 0, 0, fmt.Errorf("error during opening json work dir :%+v", filesErr)
	}
	defer files.Close()
	fileNames, fileNameErr := files.Readdirnames(10_000_000_000) // and no this one is not gonna create an array with cap string 10M
	if fileNameErr != nil {
		return 0, 0, fmt.Errorf("error during opening json work dir :%+v", filesErr)
	}
	var (
		fileRead  filesRead = 0
		byteRead  bytesRead = 0
		fileBytes           = make([]byte, 50)
		file      *os.File
		err       error
		n         int
	)
	for _, name := range fileNames {
		if len(sigCh)==1{break}
		file, err = os.Open(name)
		if err != nil {
			log.Printf("error during opening file %v with value %+v", name, err)
			file.Close()
			continue
		}
		n, err = io.ReadFull(file, fileBytes)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			log.Printf("error during reading file %v with value %+v", name, err)
			file.Close()
			continue
		}
		fileRead += 1
		byteRead += bytesRead(n)
		file.Close()
		js := &jsonwriter.Js{}
		if err = json.Unmarshal(fileBytes[:n], js); err != nil {
			log.Printf("error during unmarshal file %v with value %+v", name, err)
		}
		jsondelete.DeleteFile(name)
	}
	return fileRead, byteRead, nil
}
