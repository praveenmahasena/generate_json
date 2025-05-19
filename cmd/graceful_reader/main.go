package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"

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
	dir, dirErr := os.Open("./json")
	if dirErr != nil {
		return 0, 0, fmt.Errorf("error during opening up %v with value %+v", "./json", dirErr)
	}
	defer dir.Close()
	var (
		fileRead  uint
		bytesRead uint
	)
	for {
		if len(sigCh) == 1 {
			break
		}
		dirNames, dirNamesErr := dir.Readdirnames(10)
		if dirNamesErr != nil && !errors.Is(dirNamesErr, io.EOF) {
			break
		}
		for _, jDir := range dirNames {
			f, b, err := readSubDir(jDir, sigCh)
			if err != nil {
				log.Println(err)
			}
			fileRead += f
			bytesRead += b
		}

		if len(dirNames) < 10 {
			break
		}
	}
	deleteAll("./json")
	return fileRead, bytesRead, nil
}

func readSubDir(jDir string, sigCh chan os.Signal) (uint, uint, error) {
	p := path.Join("./json/", jDir, "./")
	dir, dirErr := os.Open(p)
	if dirErr != nil {
		return 0, 0, fmt.Errorf("error during opening dir %v with value %+v", p, dirErr)
	}
	defer dir.Close()
	var (
		fileRead  uint
		bytesRead uint
	)
	for {
		if len(sigCh) == 1 {
			break
		}
		dirNames, dirNamesErr := dir.Readdirnames(10)
		if dirNamesErr != nil && !errors.Is(dirNamesErr, io.EOF) {
			break
		}
		for _, jDir := range dirNames {
			fd := p + "/" + jDir
			b, err := os.ReadFile(fd)
			if err != nil {
				log.Printf("error during reading out file %v with value %+v", fd, err)
				deleteAll(fd)
				continue
			}
			deleteAll(fd)
			fileRead += 1
			bytesRead += uint(len(b))
			json.Unmarshal(b,&internal.Js{})
		}
		if len(dirNames) < 10 {
			break
		}
	}
	deleteAll(p)
	return fileRead, bytesRead, nil
}

func deleteAll(fName string) error {
	if err := os.RemoveAll(fName); err != nil {
		return fmt.Errorf("error during deleting file %+v", fName)
	}
	return nil
}
