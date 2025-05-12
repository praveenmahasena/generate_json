package jsonreader

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"

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
		return 0, 0, fmt.Errorf("error during enter work dir %+v", wdErr)
	}
	files, err := os.ReadDir(".")
	if err != nil {
		return 0, 0, fmt.Errorf("error during getting work dir %+v", wdErr)
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
	}
	return fr, br, nil
}

func GracefulRead(sigCh chan os.Signal, l *slog.Logger) (fileRead, byteRead, error) {
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return 0,0,fmt.Errorf("error during getting work dir %+v", wdErr)
	}
	path := path.Join(wd, "/json/")
	if err := os.Chdir(path); err != nil {
		return 0,0,fmt.Errorf("error during enter work dir %+v", wdErr)
	}
	files, err := os.ReadDir(".")
	if err != nil {
		return 0,0,fmt.Errorf("error during getting work dir %+v", wdErr)
	}
	fr := fileRead(0)
	br := byteRead(0)
	for _, f := range files {
		if len(sigCh)>=1{
			fmt.Println("cancelling...")
			break
		}
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
	}
	fmt.Println("almost cancelled....")
	return fr,br,nil
}
