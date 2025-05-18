package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/praveenmahasena/generate_json/internal"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	fileRead, byteRead, err := gracefulRead(sigCh)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	fmt.Printf("total files read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", byteRead)
}

func gracefulRead(sigCh chan os.Signal) (uint, uint, error) {
	var (
		fileRead  uint
		readBytes uint
	)
	err := filepath.Walk("./json", func(path string, info fs.FileInfo, err error) error {
		if len(sigCh)==1 {return fmt.Errorf("cancel during readin files")}
		if !info.IsDir() {
			b, err := os.ReadFile("./" + path)
			if err != nil && !errors.Is(err, io.EOF) {
				log.Printf("error during reading file %v with value %+v", path, err)
				return nil
			}
			fileRead += 1
			readBytes += uint(len(b))
			json.Unmarshal(b,&internal.Js{}) // I'm doing like this cuz we don't do much here with json data
		}
		return nil
	})
	if err!=nil{
		return fileRead, readBytes, err
	}
	deleteFile("./json")
	return fileRead, readBytes, nil
}

func deleteFile(fName string) error {
	if err := os.RemoveAll(fName); err != nil {
		return fmt.Errorf("error during deleting file %+v", fName)
	}
	return nil
}
