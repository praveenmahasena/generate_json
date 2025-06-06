// Package main ...
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

type state struct {
	path      string
	filesRead *atomic.Uint64
	bytesRead *atomic.Uint64
}

func newState(path string) *state {
	return &state{path, &atomic.Uint64{}, &atomic.Uint64{}}
}

func main() {
	path, pathErr := getPath()
	if pathErr != nil {
		fmt.Fprintln(os.Stdout, pathErr)
		syscall.Exit(-1)
	}
	s := newState(path)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(cancel context.CancelFunc) {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel()
	}(cancel)
	nameCh := make(chan string)
	doneCh := make(chan bool)
	defer close(doneCh)
	go s.ProcessAndDelete(nameCh, doneCh)
	if err := s.Read(ctx, nameCh); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	close(nameCh)
	<-doneCh
	s.ShowStats()
}

func (s *state) Read(ctx context.Context, nameCh chan<- string) error {
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
		prossesDirectories(ctx, subDirNames, nameCh)
	}
	return nil
}

func (s *state) ProcessAndDelete(fileNameCh <-chan string, done chan bool) {
	for fileName := range fileNameCh {
		b, err := os.ReadFile(fileName)
		if err != nil {
			log.Printf("error during reading file %v with value %v", fileName, err)
			continue
		}
		if err := json.Unmarshal(b, &internal.Js{}); err != nil {
			log.Printf("error during wring bytes file %v with value %v", fileName, err)
			continue
		}
		s.filesRead.Add(1)
		s.bytesRead.Add(uint64(len(b)))
		if err := os.Remove(fileName); err != nil {
			log.Printf("error during deleting file %v with value %#v", fileName, err)
		}
	}
	done <- true
}

func (s *state) ShowStats() {
	fmt.Printf("\ntotal file read %v\n", s.filesRead.Load())
	fmt.Printf("total bytes read %v \n", s.bytesRead.Load())
}

func prossesDirectories(ctx context.Context, subDirNames []string, nameCh chan<- string) {
	for _, subDirName := range subDirNames {
		if ctx.Err() != nil {
			break
		}
		prossesDirectory(ctx, subDirName, nameCh)
	}
}

func prossesDirectory(ctx context.Context, subDirName string, nameCh chan<- string) {
	p := "./json/" + subDirName
	subDirectory, subDirectoryErr := os.Open(p)
	if subDirectoryErr != nil {
		log.Printf("error during opening dir :%v with value %+v", p, subDirectoryErr)
		return
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
		processFileNames(ctx, p, fileNames, nameCh)
	}
}

func processFileNames(ctx context.Context, p string, fileNames []string, fileNameCh chan<- string) {
	for _, fileName := range fileNames {
		if ctx.Err() != nil {
			break
		}
		fileNameCh <- path.Join(p, "/", fileName)
	}

}

func getPath() (string, error) {
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return "", fmt.Errorf("error during getting working dir with value %#v", wdErr)
	}
	return wd + "/json", nil
}
