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
		fileRead  uint
		bytesRead uint
		fileskip  uint
	)

	for {
		if len(sigCh) == 1 {
			l.Error("closing off syscall.SIGINT")
			break
		}
		subDirNames, subDirNameErr := jsonDir.Readdirnames(10)
		if subDirNameErr != nil {
			if errors.Is(subDirNameErr, io.EOF) {
				l.Info("all files processed")
				break
			}
			fileskip += 1
			l.Error("error during reading ./json sub dirs", "error value", subDirNameErr, "process", "skipping...")
			continue
		}
	inner:
		for _, subDirs := range subDirNames {
			if len(sigCh) == 1 {
				break inner
			}
			p := "./json/" + subDirs
			subDirNames, subDirNamesErr := os.Open(p) // I should take care of (*os.File).Close() here but later not now
			if subDirNamesErr != nil {
				fileskip += 1
				l.Error("error during opening dir", "dir name", p, "error value", subDirNamesErr, "process", "skipping...")
				continue inner
			}
		deeper:
			for {
				if len(sigCh) == 1 {
					break deeper
				}
				fileCollection, fileCollectionErr := subDirNames.Readdirnames(10)
				if fileCollectionErr != nil {
					if errors.Is(fileCollectionErr, io.EOF) {
						break deeper
					}
					fileskip += 1
					l.Error("error during file name read", "dir name", p, "error value", fileCollectionErr, "process", "skipping...")
					continue deeper
				}
			fileLoop:
				for _, fileName := range fileCollection {
					if len(sigCh) == 1 {
						break fileLoop
					}
					b, err := os.ReadFile(p + "/" + fileName)
					if err != nil {
						l.Error("error during reading file", "file name", p+"/"+fileName, "error value", err, "process", "skipping...")
						fileskip += 1
						continue fileLoop
					}
					if err := json.Unmarshal(b, &internal.Js{}); err != nil {
						fileskip += 1
						l.Error("error during marshelling file", "file name", p+"/"+fileName, "error value", err, "process", "skipping...")
						continue fileLoop
					}
					fileRead += 1
					bytesRead += uint(len(b))
				}
			}
		}
	}
	if fileskip >= 1 {
		if err := os.RemoveAll("./json"); err != nil {
			return fileRead, bytesRead, fmt.Errorf("error during removing json dir %+v", err)
		}
	}
	return fileRead, bytesRead, nil
}

// And that you have a for loop calling multiple functions to pricess a file that can all return errors is making your code more difficult to reason and maintain.
// as your message said I changes it to more simple way of doing this

//5) only delete the file if all processing of the filename has no errors.  For now since it’s confusing the issue don’t worry about removing the subdirectory if it’s empty.

// This is my point about creating a single function that does that processing so you can change the processing and still check for a single error.
// I assume you meant delete the dir if theres all file read so I am checking if any file went wrong if so I do not delete the root dir
