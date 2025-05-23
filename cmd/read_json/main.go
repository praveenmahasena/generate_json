package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
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
		fileRead  *atomic.Uint64
		bytesRead *atomic.Uint64
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
		if err := prossesDirectories(subDirNames, fileRead, bytesRead, l); err != nil {
			// I did not change any kind of logic at all
			// created prossesDirectories([]string, *atomic.Uint64,*slog.Logger) error
			// passed in atomic pointers
			// according to the previous logic we handle errors just by logging and not returning
			// please give me some suggestions
		}
	}
	return uint(fileRead.Load()), uint(bytesRead.Load()), nil
}

func prossesDirectories(subDirNames []string, fileRead, bytesRead *atomic.Uint64, l *slog.Logger) error {
	for _, subDirs := range subDirNames {
		p := "./json/" + subDirs
		subDirNames, subDirNamesErr := os.Open(p) // I should take care of (*os.File).Close() here but later not now
		if subDirNamesErr != nil {
			l.Error("error during opening dir", "dir name", p, "error value", subDirNamesErr, "process", "skipping...")
			continue
			// what are we doing here with this "subDirNamesErr"?
			// we are skipping one subdirs since we do not get to open it no need to return here
			// it would messup all the other dirs that should be read
		}
	deeper:
		for {
			fileCollection, fileCollectionErr := subDirNames.Readdirnames(10)
			if fileCollectionErr != nil {
				if errors.Is(fileCollectionErr, io.EOF) {
					break deeper
					// io.EOF error does not matter
				}
				l.Error("error during file name read", "dir name", p, "error value", fileCollectionErr, "process", "skipping...")
				continue deeper
				// what are we doing here with this "fileCollectionErr"?
				// we are skipping one of 10 subdirs since we do not get to open it no need to return here
				// it would messup all the other dirs that should be read
			}
		fileLoop:
			for _, fileName := range fileCollection {
				b, err := os.ReadFile(p + "/" + fileName)
				if err != nil {
					l.Error("error during reading file", "file name", p+"/"+fileName, "error value", err, "process", "skipping...")
					continue fileLoop
				}
				if err := json.Unmarshal(b, &internal.Js{}); err != nil {
					l.Error("error during marshelling file", "file name", p+"/"+fileName, "error value", err, "process", "skipping...")
					continue fileLoop
				}
				fileRead.Add(1)
				bytesRead.Add(uint64(len(b)))
			}
		}
		subDirNames.Close()
	}
	return nil
}
