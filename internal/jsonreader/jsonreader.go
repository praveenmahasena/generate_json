package jsonreader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/praveenmahasena/generate_json/internal/jsondelete"
	"github.com/praveenmahasena/generate_json/internal/jsonwriter"
)

type (
	fileRead uint
	byteRead uint
)

// TODO: refactor get similar func to util package

// I don't wanna do named return params
// it's just not my style if it's nessocery please let me know
// the 1st one is amount of files read
// the 2st one is amount of file bytes read
// the 3st one is any kind of error
// I have created types for it on top for this

func Read(l *slog.Logger) (fileRead, byteRead, error) {
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return 0, 0, fmt.Errorf("error during getting work dir %+v", wdErr)
	}
	path := path.Join(wd, "/json/")
	if err := os.Chdir(path); err != nil {
		return 0, 0, fmt.Errorf("error during enter work dir %+v", err)
	}
	files, err := os.ReadDir(".")
	if err != nil {
		return 0, 0, fmt.Errorf("error during getting work dir %+v", err)
	}
	fr := fileRead(0)
	br := byteRead(0)
	for _, f := range files {
		fr += 1
		b, err := os.ReadFile(f.Name())
		if err != nil {
			errStr := fmt.Errorf("error during reading file %+v", err)
			// I'll be doing logging here for error handling since I do not see a purpose of bubbling up
			l.Error("error", errStr.Error(), "")
			continue
		}
		br += byteRead(len(b))
		js := jsonwriter.Js{}
		if err := json.Unmarshal(b, &js); err != nil {
			// this error handling is not nessocery
			// but still just gonna log into stdErr if json goes wrong that's all
			errStr := fmt.Errorf("error during unmarshalling file %+v with value %+v", f, err)
			l.Error("error", errStr.Error(), "")
		}
		if err := jsondelete.DeleteFile(f.Name()); err != nil {
			errStr := fmt.Errorf("error during deleting file %+v with value %+v", f.Name(), err)
			l.Error("error", errStr.Error(), "")
		}
	}
	return fr, br, nil
}

// fileRead, byteRead, error
func GracefulRead(sigCh chan os.Signal, l *slog.Logger) (uint, uint, error) {
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return 0, 0, fmt.Errorf("error during getting work dir %+v", wdErr)
	}
	p := path.Join(wd, "/json/")
	if err := os.Chdir(p); err != nil {
		return 0, 0, fmt.Errorf("error during moving into json work dir %+v", err)
	}
	var (
		errCount      uint8
		fileTrack     uint
		fName         string
		file          *os.File
		err           error
		buf           []byte = make([]byte, 100)
		n             int
		fileBytesRead uint
		fileRead      uint
	)

	for {
		if len(sigCh) == 1 {
			fmt.Println("cancelling...")
			break
		}
		if errCount >= 1 {
			break
		}
		fName, fileTrack = fmt.Sprintf("%v.json", fileTrack), fileTrack+1
		file, err = os.OpenFile(fName, os.O_RDONLY, 0)
		if err != nil {
			l.Error("skipping file", fName, "")
			err = nil
			errCount += 1
			continue
		}
		n, err = io.ReadFull(file, buf)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			l.Error("skipping file", fName, "")
			err = nil
			continue
		}
		fileBytesRead += uint(n)
		fileRead += 1
		file.Close()
		js := &jsonwriter.Js{}

		err = json.Unmarshal(buf[:n], js)
		if err != nil {
			l.Error("error during unmarshal json", err.Error(), "")
			err = nil
		}
		if err = jsondelete.DeleteFile(fName); err != nil {
			l.Error("error during deleting json file", err.Error(), "")
			err = nil
		}
	}
	return fileRead, fileBytesRead, nil
}
