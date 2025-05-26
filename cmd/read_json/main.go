package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"sync/atomic"

	"github.com/praveenmahasena/generate_json/internal"
)

func main() {
	l := internal.NewLogger(os.Stdout, true, 1)
	fileRead, bytesRead, err := read(l)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Printf("total file read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", bytesRead)
}

func read(l *slog.Logger) (uint, uint, error) {
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
		subDirNames, subDirNameErr := jsonDir.Readdirnames(10)
		if subDirNameErr != nil {
			if errors.Is(subDirNameErr, io.EOF) {
				l.Info("all files processed")
				break
			}
			l.Error("error during reading ./json sub dirs", "error value", subDirNameErr, "process", "skipping...")
			continue
		}
		prossesDirectories(subDirNames, &fileRead, &bytesRead, l)
	}
	return uint(fileRead.Load()), uint(bytesRead.Load()), nil
}

func prossesDirectories(subDirNames []string, fileRead, bytesRead *atomic.Uint64, l *slog.Logger) error {
	for _, subDirName := range subDirNames {
		p := "./json/" + subDirName
		subDirectory, subDirectoryErr := os.Open(p)
		if subDirectoryErr != nil {
			l.Error("error during opening", "error value", subDirectoryErr, "process", "skipping...")
			continue
		}
		prossesDirectory(p, subDirectory, fileRead, bytesRead, l)
		subDirectory.Close()
	}
	return nil
}

func prossesDirectory(p string, subDirectories *os.File, fileRead, bytesRead *atomic.Uint64, l *slog.Logger) error {
	defer subDirectories.Close()
	for {
		fileNames, fileNamesErr := subDirectories.Readdirnames(10)
		if fileNamesErr != nil {
			if errors.Is(fileNamesErr, io.EOF) {
				break
			}
			l.Error("error during getting file names", "error value", fileNamesErr, "process", "skipping...")
			continue
		}
		processFileNames(p, fileNames, fileRead, bytesRead)
	}
	return nil
}

func processFileNames(p string, fileNames []string, fileRead, bytesRead *atomic.Uint64) error {
	for _, fileName := range fileNames {
		fileName = path.Join(p, "/", fileName)
		if err := processDeleteFile(fileName, fileRead, bytesRead); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func processDeleteFile(fileName string, fileRead, bytesRead *atomic.Uint64) error {
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
		return fmt.Errorf("error during deleting file %v with error value %+v", fileName, err)
	}
	return nil
}
