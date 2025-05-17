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
	dirRead, fileRead, byteRead, err := gracefulRead(sigCh)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	fmt.Printf("total dir read %v \n", dirRead)
	fmt.Printf("total files read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", byteRead)
}

func gracefulRead(sigCh chan os.Signal) (uint, uint, uint, error) {
	if err := os.Chdir("./json"); err != nil {
		return 0, 0, 0, fmt.Errorf("cannot move into json dir %+v", err)
	}
	file, fileErr := os.Open(".")
	if fileErr != nil {
		return 0, 0, 0, fmt.Errorf("cannot open dir json %+v", fileErr)
	}
	defer file.Close()
	dirNames, dirNameErr := file.Readdirnames(100)
	if dirNameErr != nil {
		return 0, 0, 0, fmt.Errorf("cannot read json dir names %+v", dirNameErr)
	}
	var (
		dirRead  uint
		byteRead uint
		fileRead uint
	)
	for _, dir := range dirNames {
		if len(sigCh)==1{break}
		f, fErr := os.Open("./" + dir)
		if fErr != nil {
			log.Printf("error during opening in dir %v with value %+v", dir, fErr)
			continue
		}
		fileNames, err := f.Readdirnames(10_000)
		if err != nil {
			log.Printf("error during reading all json file name in dir %v with value %+v", dir, err)
			f.Close()
			continue
		}
		br, fr := readDirFiles(fileNames, dir,sigCh)
		f.Close()
		byteRead += br
		fileRead += fr
		dirRead += 1
		deleteFile(dir)
	}
	return dirRead, fileRead, byteRead, nil
}

func readDirFiles(fileNames []string, dirName string,sigCh chan os.Signal) (uint, uint) {
	if err := os.Chdir(dirName); err != nil {
		log.Printf("error during moving into %v with %+v", dirName, err)
		return 0, 0
	}
	defer os.Chdir("../")
	buf := make([]byte, 50)
	var fr uint
	var br uint
	for _, fn := range fileNames {
		if len(sigCh)==1{break}
		f, e := os.OpenFile(fn,os.O_RDONLY,0666)
		if e != nil {
			log.Printf("error during opening up file %v in dir %v with %+v", dirName, fn, e)
			continue
		}
		n, err := io.ReadFull(f, buf)
		f.Close()
		if err != nil && !errors.Is(err,io.ErrUnexpectedEOF) {
//			log.Printf("error during reading file %v in dir %v with %+v", dirName, fn, err)
		}
		fr += 1
		br += uint(n)
		json.Unmarshal(buf[:n],&internal.Js{})
	}
	return br, fr
}
func deleteFile(fName string) error {
	if err := syscall.Unlink(fName); err != nil {
		return fmt.Errorf("error during deleting file %+v", fName)
	}
	return nil
}
