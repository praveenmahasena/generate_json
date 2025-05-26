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
	fileRead := atomic.Uint64{}
	bytesRead := atomic.Uint64{}
	sigCh := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // let's face it we do not need another leak

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func(cancel context.CancelFunc, sigCh chan os.Signal) {
		<-sigCh
		fmt.Println("cancel signal has been sent")
		cancel()
	}(cancel, sigCh)

	if err := gracefulRead(ctx, &fileRead, &bytesRead); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Printf("file read %v \n", fileRead.Load())
	fmt.Printf("bytes read %v \n", bytesRead.Load())
}

func gracefulRead(ctx context.Context, fileRead, bytesRead *atomic.Uint64) chan error {
	errCh := make(chan error, 1)
	jsonDir, jsonDirErr := os.Open("./json")
	if jsonDirErr != nil {
		errCh <- fmt.Errorf("error during opeing ./json dir with value %+v", jsonDirErr)
		return errCh
	}
	defer jsonDir.Close()

	for {
		<-ctx.Done()
		subDirNames, subDirNameErr := jsonDir.Readdirnames(10)
		if subDirNameErr != nil {
			if errors.Is(subDirNameErr, io.EOF) {
				log.Println("all files are processed")
				break
			}
			log.Printf("error during reading subdir in ./json with error value %+v", subDirNameErr)
			continue
		}
		prossesDirectories(ctx,subDirNames, fileRead, bytesRead )
	}
	return nil
}

func prossesDirectories(ctx context.Context,subDirNames []string, fileRead, bytesRead *atomic.Uint64) error {
	for _, subDirName := range subDirNames {
		if bool.Load() {
			log.Println("cancelling due to syscall.SIGINT signal")
			break
		} // we are doing bubble up here
		if err := prossesDirectory(subDirName, fileRead, bytesRead, bool); err != nil {
			log.Panicln(err)
		}
	}
	return nil
}

func prossesDirectory(subDirName string, fileRead, bytesRead *atomic.Uint64, bool *atomic.Bool) error {
	p := "./json/" + subDirName
	subDirectory, subDirectoryErr := os.Open(p)
	if subDirectoryErr != nil {
		return fmt.Errorf("error during opening dir :%v with value %+v", p, subDirectoryErr)
	}
	defer subDirectory.Close()
	for {
		if bool.Load() {
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
		processFileNames(p, fileNames, fileRead, bytesRead, bool)
	}
	return nil
}

func processFileNames(p string, fileNames []string, fileRead, bytesRead *atomic.Uint64, bool *atomic.Bool) error {
	for _, fileName := range fileNames {
		if bool.Load() {
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
