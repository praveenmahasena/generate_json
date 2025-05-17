package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
)

type (
	filesRead uint
	bytesRead uint
)

type Js struct {
	ID uint `json:"id"`
}

func main() {
	fileRead, byteRead, err := read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// there is no a return statement here it's not an error
		// I just wanna show it as the file read are 0
	}
	fmt.Printf("total files read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", byteRead)
}

func read() (filesRead, bytesRead, error) {
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
		js := &Js{}
		if err = json.Unmarshal(fileBytes[:n], js); err != nil {
			log.Printf("error during unmarshal file %v with value %+v", name, err)
		}
		deleteFile(name)
	}
	return fileRead, byteRead, nil
}

func deleteFile(fn string) error {
	if err := syscall.Unlink(fn); err != nil {
		return fmt.Errorf("error during deleting file %v with value %v", fn, err)
	}
	return nil
}
