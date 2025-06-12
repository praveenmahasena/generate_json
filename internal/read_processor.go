package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync/atomic"
)

type DirectoryReaderInterface interface {
	Read(context.Context, chan<- string) error
	ShowStats()
	AddState(uint64, uint64)
}

func (s *state) ShowStats() {
	fmt.Printf("\n fileRead %v", s.fileRead.Load())
	fmt.Printf("\n bytesRead %v", s.bytesRead.Load())
}

type state struct {
	fileRead, bytesRead *atomic.Uint64
}

type DirectoryReader struct {
	path string
	*state
}

func (s *state) AddState(fileRead, bytesRead uint64) {
	s.fileRead.Add(fileRead)
	s.bytesRead.Add(bytesRead)
}

func (s *state) GetStats() {
	fmt.Printf("\n file read %v", s.fileRead)
	fmt.Printf("\n bytes read %v", s.bytesRead)
}

func NewdirectoryReader(path string, fileRead, bytesRead *atomic.Uint64) *DirectoryReader {
	return &DirectoryReader{path, &state{fileRead, bytesRead}}
}

func (d *DirectoryReader) Read(ctx context.Context, nameCh chan<- string) error {
	jsonDir, jsonDirErr := os.Open("./json")
	if jsonDirErr != nil {
		return fmt.Errorf("error during opeing ./json dir with value %+v", jsonDirErr)
	}
	defer jsonDir.Close()

	for {
		if ctx.Err() == context.Canceled {
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

func prossesDirectories(ctx context.Context, subDirNames []string, nameCh chan<- string) {
	for _, subDirName := range subDirNames {
		if ctx.Err() == context.Canceled {
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
		if ctx.Err() == context.Canceled {
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
		if ctx.Err() == context.Canceled {
			break
		}
		fileNameCh <- path.Join(p, "/", fileName)
	}
}
