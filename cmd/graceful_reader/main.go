package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/praveenmahasena/generate_json/internal"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	fileRead, byteRead, err := gracefulRead(sigCh)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	fmt.Println(fileRead, byteRead)
}

func gracefulRead(sigCh chan os.Signal) (uint, uint, error) {
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
		fileRead  uint = 0
		byteRead  uint = 0
		fileBytes      = make([]byte, 50)
		file      *os.File
		err       error
		n         int
	)
	for _, name := range fileNames {
		if len(sigCh) == 1 {
			break
		}
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
		byteRead += uint(n)
		file.Close()
		js := &internal.Js{}
		if err = json.Unmarshal(fileBytes[:n], js); err != nil {
			log.Printf("error during unmarshal file %v with value %+v", name, err)
		}
		deleteFile(name)
	}
	return fileRead, byteRead, nil
}

func deleteFile(fName string) error {
	if err := syscall.Unlink(fName); err != nil {
		return fmt.Errorf("error during deleting file %+v",fName)
	}
	return nil
}
