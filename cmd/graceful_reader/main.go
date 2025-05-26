package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"sync/atomic"

	"github.com/praveenmahasena/generate_json/internal"
)

func main() {
	l := internal.NewLogger(os.Stdout, true, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	fileRead, bytesRead, err := gracefulRead(sigCh, l)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Printf("total file read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", bytesRead)
}

func gracefulRead(sigCh chan os.Signal, l *slog.Logger) (uint, uint, error) {
	jsonDir, jsonDirErr := os.Open("./json")
	if jsonDirErr != nil {
		return 0, 0, fmt.Errorf("error during opeing ./json dir with value %+v", jsonDirErr)
	}
	defer jsonDir.Close()

	var (
		fileRead  = atomic.Uint64{}
		bytesRead = atomic.Uint64{}
	)

	for {
		if len(sigCh) == 1 {
			break
		}
		subDirNames, subDirNameErr := jsonDir.Readdirnames(10)
		if subDirNameErr != nil {
			if errors.Is(subDirNameErr, io.EOF) {
				l.Info("all files processed")
				break
			}
			l.Error("error during reading ./json sub dirs", "error value", subDirNameErr, "process", "skipping...")
			continue
		}
		prossesDirectories(subDirNames, &fileRead, &bytesRead, sigCh, l)
	}
	return uint(fileRead.Load()), uint(bytesRead.Load()), nil
}

func prossesDirectories(subDirNames []string, fileRead, bytesRead *atomic.Uint64, sigCh chan os.Signal, l *slog.Logger) error {
	for _, subDirName := range subDirNames {
		if len(sigCh) == 1 {
			l.Info("cancelling due to syscall.SIGINT signal")
			break
		} // we are doing bubble up here
		if err := prossesDirectory(subDirName, fileRead, bytesRead, sigCh, l); err != nil {
			log.Panicln(err)
		}
	}
	return nil
}

func prossesDirectory(subDirName string, fileRead, bytesRead *atomic.Uint64, sigCh chan os.Signal, l *slog.Logger) error {
	p := "./json/" + subDirName
	subDirectory, subDirectoryErr := os.Open(p)
	if subDirectoryErr != nil {
		return fmt.Errorf("error during opening dir :%v with value %+v", p, subDirectoryErr)
	}
	defer subDirectory.Close()
	for {
		if len(sigCh) == 1 {
			break
		}
		fileNames, fileNamesErr := subDirectory.Readdirnames(10)
		if fileNamesErr != nil {
			if errors.Is(fileNamesErr, io.EOF) {
				break
			}
			l.Error("error during getting file names", "error value", fileNamesErr, "process", "skipping...")
			continue
		}
		processFileNames(p, fileNames, fileRead, bytesRead, sigCh)
	}
	return nil
}

func processFileNames(p string, fileNames []string, fileRead, bytesRead *atomic.Uint64, sigCh chan os.Signal) error {
	for _, fileName := range fileNames {
		if len(sigCh) == 1 {
			break
		}
		fileName = path.Join(p, "/", fileName)
		if err := processAndRemoveFile(fileName, fileRead, bytesRead); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func processAndRemoveFile(fileName string, fileRead, bytesRead *atomic.Uint64) error {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error during reading file %v with value %+v", fileName, err)
	}
	if err := json.Unmarshal(b, &internal.Js{}); err != nil {
		return fmt.Errorf("error during Unmarshal file %v with value %+v", fileName, err)
	}
	fileRead.Add(1)
	bytesRead.Add(uint64(len(b)))
	if err := os.Remove(fileName); err != nil {
		return fmt.Errorf("error during deleting off file %v with value %+v", fileName, err)
	}
	return nil
}
