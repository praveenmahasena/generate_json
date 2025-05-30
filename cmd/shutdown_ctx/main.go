package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"sync/atomic"
	"syscall"

	"github.com/praveenmahasena/generate_json/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // let's not have a context leak shall we?

	go func(cancel context.CancelFunc) {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("cancelling....")
		cancel() // or shall I do defer and make it in the top of func? // for now it does not matter for me
	}(cancel)

	var fileRead, bytesRead atomic.Uint64
	if err := gracefulRead(ctx, &fileRead, &bytesRead); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("total file read %v \n", fileRead.Load())
	fmt.Printf("total bytes read %v \n", bytesRead.Load())
}

func gracefulRead(ctx context.Context, fileRead, bytesRead *atomic.Uint64) error {
	jsonDir, jsonDirErr := os.Open("./json")
	if jsonDirErr != nil {
		return fmt.Errorf("error during opeing ./json dir with value %+v", jsonDirErr)
	}
	defer jsonDir.Close()

	for {
		if ctx.Err() != nil {
			break
		}
		subDirNames, subDirNameErr := jsonDir.Readdirnames(10)
		if subDirNameErr != nil {
			if errors.Is(subDirNameErr, io.EOF) {
				log.Println("all files are processed")
				break
			}
			log.Printf("error during reading subdir in ./json with error value %+v", subDirNameErr)
			continue
		}
		prossesDirectories(ctx, subDirNames, fileRead, bytesRead)
	}
	return nil
}

func prossesDirectories(ctx context.Context, subDirNames []string, fileRead, bytesRead *atomic.Uint64) error {
	for _, subDirName := range subDirNames {
		if ctx.Err() != nil {
			break
		}
		if err := prossesDirectory(ctx, subDirName, fileRead, bytesRead); err != nil {
			log.Panicln(err)
		}
	}
	return nil
}

func prossesDirectory(ctx context.Context, subDirName string, fileRead, bytesRead *atomic.Uint64) error {
	p := "./json/" + subDirName
	subDirectory, subDirectoryErr := os.Open(p)
	if subDirectoryErr != nil {
		return fmt.Errorf("error during opening dir :%v with value %+v", p, subDirectoryErr)
	}
	defer subDirectory.Close()
	for {
		if ctx.Err() != nil {
			break
		}
		fileNames, fileNamesErr := subDirectory.Readdirnames(10)
		if fileNamesErr != nil {
			if errors.Is(fileNamesErr, io.EOF) {
				break
			}
			log.Printf("error during getting file names with error value %+v process skipping...", fileNamesErr)
			continue
		}
		processFileNames(ctx, p, fileNames, fileRead, bytesRead)
	}
	return nil
}

func processFileNames(ctx context.Context, p string, fileNames []string, fileRead, bytesRead *atomic.Uint64) error {
	for _, fileName := range fileNames {
		if ctx.Err() != nil {
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
