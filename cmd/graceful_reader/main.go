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
	var dirCount uint
	for {
		if len(sigCh) == 1 {
			l.Info("cancelling due to signal")
			break
		}
		dirNames, dirNamesErr := dir.Readdirnames(10)
		if dirNamesErr != nil && !errors.Is(dirNamesErr, io.EOF) {
			return fileRead, bytesRead, fmt.Errorf("error during reading out directories with value %+v", dirNamesErr)
		}
		for _, jDir := range dirNames {
			f, b, err := readSubDir(jDir, sigCh,l)
			if err != nil {
				log.Println(err)
			}
			fileRead += f
			bytesRead += b
		}

		if len(dirNames) < 10 {
			dirCount += 1
			l.Info("sucessfully read a dir", "dir amount read", dirCount)
			break
		}
	}
	// I am not a if else if fan but I'm okay doing it here cuz this one has very small logic block
	if d, err := dir.Readdirnames(10); err != nil {
		return fileRead, bytesRead, fmt.Errorf("error during deleting process of checking on json dir %+v", err)
	} else if len(d) > 0 {
		return fileRead, bytesRead, fmt.Errorf("unread files exists in some dirs fully deletion cancelled with value %+v", err)
	}
	if err := deleteAll("./json"); err != nil {
		return fileRead, bytesRead, fmt.Errorf("error during deleting root of json dir collection with value %+v", err)
	}
	return fileRead, bytesRead, nil
}

func readSubDir(jDir string, sigCh chan os.Signal,l *slog.Logger) (uint, uint, error) {
	p := path.Join("./json/", jDir, "./")
	dir, dirErr := os.Open(p)
	if dirErr != nil {
		return 0, 0, fmt.Errorf("error during opening dir %v with value %+v", p, dirErr)
	}
	defer dir.Close()
	var (
		fileRead  uint
		bytesRead uint
	)
	for {
		if len(sigCh) == 1 {
			// no need to do any logs here since the parent loop is gonna run one more time before being cancelled so and it has the cancel log
			break
		}
		dirNames, dirNamesErr := dir.Readdirnames(10)
		if dirNamesErr != nil && !errors.Is(dirNamesErr, io.EOF) {
			break
		}
		for _, jDir := range dirNames {
			fd := p + "/" + jDir
			b, err := os.ReadFile(fd)
			if err != nil {
				log.Printf("error during reading out file %v with value %+v", fd, err)
				deleteAll(fd)
				continue
			}
			deleteAll(fd)
			fileRead += 1
			bytesRead += uint(len(b))
			json.Unmarshal(b, &internal.Js{})
		}
		if len(dirNames) < 10 {
			break
		}
	}
	deleteAll(p)
	return fileRead, bytesRead, nil
}

func deleteAll(fName string) error {
	if err := os.RemoveAll(fName); err != nil {
		return fmt.Errorf("error during deleting file %+v", fName)
	}
	return nil
}
