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

	"github.com/praveenmahasena/generate_json/internal"
)

func main() {
	l := internal.NewLogger(os.Stdout, true, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	fileRead, byteRead, err := gracefulRead(sigCh, l)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	fmt.Printf("total files read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", byteRead)
}

func gracefulRead(sigCh chan os.Signal, l *slog.Logger) (uint, uint, error) {
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
			l.Info("cancelling due to signal")
			break
		}
		dirNames, dirNamesErr := dir.Readdirnames(-1)
		if dirNamesErr != nil && !errors.Is(dirNamesErr, io.EOF) {
			return fileRead, bytesRead, fmt.Errorf("error during reading out directories with value %+v", dirNamesErr)
		}
		for _, jDir := range dirNames {
			f, b, err := readSubDir(jDir, sigCh, l)
			if err != nil {
				log.Println(err)
			}
			fileRead += f
			bytesRead += b
		}

		if len(dirNames) < 10 {
			l.Info("all dir read...")
			break
		}
	}
	// I can do somekind of bubble up value and make the "./json" dir get deleted but i do not think it's that important since all the other dirs are getting cleaned up :)
	return fileRead, bytesRead, nil
}

func readSubDir(jDir string, sigCh chan os.Signal, l *slog.Logger) (uint, uint, error) {
	p := path.Join("./json", jDir, "/")
	dir, dirErr := os.Open(p)
	if dirErr != nil {
		return 0, 0, fmt.Errorf("error during opening up %v with value %+v", p, dirErr)
	}
	defer dir.Close() // Im gonna delete the dir if the file amount I got is == 0 if I do this while having this in a defer stack would get me an error?
	// for safety reason I'll have some more close statements in places
	// yes closing a closed *os.File is gonna get me error but guess what this kind of error is meaningless
	// so I do not have to manage it
	dirNames, dirNameErr := dir.Readdirnames(-1)
	if dirNameErr != nil && !errors.Is(dirNameErr,io.EOF){
		return 0, 0, fmt.Errorf("error during getting up %v all the files with value of %+v", p, dirNameErr)
	}
	var (
		failedFiles uint
		fileRead    uint
		bytesRead   uint
	)
	for _, name := range dirNames {
		if len(sigCh) == 1 {
			break
		}
		b, err := os.ReadFile(p + "/" + name)
		if err != nil {
			failedFiles += 1
			l.Error("error during read file", "file name", p+"/"+name, "error value", err)
			continue
		}
		if err := json.Unmarshal(b, &internal.Js{}); err != nil {
			failedFiles += 1
			l.Error("error during decoding json into file", "file name", p+"/"+name, "error value", err)
			continue
		}
		fileRead += 1
		bytesRead += uint(len(b))
	}
	if failedFiles > 0 {
		return fileRead, bytesRead, fmt.Errorf("dir %v coulndnt be deleted", p)
	}
	dir.Close()
	if err := deleteAll(p); err != nil {
		return fileRead, bytesRead, fmt.Errorf("error during deleting dir %v with value %+v", p, err)
	}
	return fileRead, bytesRead, nil
}

func deleteAll(fName string) error {
	if err := os.RemoveAll(fName); err != nil {
		return fmt.Errorf("error during deleting file %+v", fName)
	}
	return nil
}
