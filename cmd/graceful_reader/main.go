package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"

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
	fmt.Printf("total file read %v", fileRead)
	fmt.Printf("total bytes read %v", bytesRead)
}

func gracefulRead(sigCh chan os.Signal, l *slog.Logger) (uint, uint, error) {
	jsonDir, jsonDirErr := os.Open("./json")
	if jsonDirErr != nil {
		return 0, 0, fmt.Errorf("error during opeing ./json dir with value %+v", jsonDirErr)
	}
	defer jsonDir.Close()

	var (
		fileRead  uint
		bytesRead uint
	)

	for {
		if len(sigCh)==1{
			l.Error("closing off syscall.SIGINT")
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
		inner:
		for _, subDirs := range subDirNames {
			if len(sigCh)==1{
				break inner
			}
			p := "./json/" + subDirs
			subDirNames, subDirNamesErr := os.Open(p)
			if subDirNamesErr != nil {
				l.Error("error during opening dir", "dir name", p, "error value", subDirNamesErr, "process", "skipping...")
				continue inner
			}
			deeper:
			for {
				if len(sigCh)==1{
					break deeper
				}
				fileCollection, fileCollectionErr := subDirNames.Readdirnames(10)
				if fileCollectionErr != nil {
					if errors.Is(fileCollectionErr, io.EOF) {
						break deeper
					}
					l.Error("error during file name read", "dir name", p, "error value", fileCollectionErr, "process", "skipping...")
					continue deeper
				}
				fileLoop:
				for _, fileName := range fileCollection {
					if len(sigCh)==1{
						break fileLoop
					}
					b, err := os.ReadFile(p + "/" + fileName)
					if err != nil {
						l.Error("error during reading file", "file name", p+"/"+fileName, "error value", err, "process", "skipping...")
						continue fileLoop
					}
					if err := json.Unmarshal(b, &internal.Js{}); err != nil {
						l.Error("error during marshelling file", "file name", p+"/"+fileName, "error value", err, "process", "skipping...")
						continue fileLoop
					}
					fileRead += 1
					bytesRead += uint(len(b))
				}
			}
		}
	}
	return fileRead, bytesRead, nil
}
