package jsonreader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"

	"github.com/praveenmahasena/generate_json/internal/jsondelete"
	"github.com/praveenmahasena/generate_json/internal/jsonwriter"
)

type (
	fileRead uint
	byteRead uint
)

// these approach is more like writing the ls cmd out put into a file and reading it in streaming way one by one which makes it more performance and less memory usage
// but the file does not refresh or change state so I had to close and reopen it please help
func Read(l *slog.Logger) (fileRead, byteRead, error) {

	temp, tempErr := os.OpenFile("temp", os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0660)
	if tempErr != nil {
		return 0, 0, fmt.Errorf("error during opening up temp file %+v", tempErr)
	}
	cmd := exec.Command("ls", "-1", "./json/")
	cmd.Stdout = temp

	if err := cmd.Run(); err != nil {
		return 0, 0, fmt.Errorf("error during running the list cmd %+v", err)
	}
	temp.Close()
	t, tErr := os.Open("temp")
	if tErr != nil {
		return 0, 0, fmt.Errorf("error during Opening the list file %+v", tErr)
	}
	r := bufio.NewScanner(t)
	r.Split(bufio.ScanLines)

	var (
		f         *os.File
		err       error
		filesRead fileRead
		bytesRead byteRead
	)
	b := bytes.NewBuffer(nil)
	for r.Scan() {
		f, err = os.Open("./json/" + r.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}
		io.Copy(b, f)
		f.Close()
		filesRead += 1
		bytesRead += byteRead(b.Len())
		js := &jsonwriter.Js{}
		json.Unmarshal(b.Bytes(),js)
		jsondelete.DeleteFile("./json/"+r.Text())
	}
	t.Close()
	jsondelete.DeleteFile("temp")
	return filesRead, bytesRead, err
}

// fileRead, byteRead, error
// read the read func comment please
func GracefulRead(sigCh chan os.Signal, l *slog.Logger) (fileRead, byteRead, error) {
	temp, tempErr := os.OpenFile("temp", os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0660)
	if tempErr != nil {
		return 0, 0, fmt.Errorf("error during opening up temp file %+v", tempErr)
	}
	cmd := exec.Command("ls", "-1", "./json/")
	cmd.Stdout = temp

	if err := cmd.Run(); err != nil {
		return 0, 0, fmt.Errorf("error during running the list cmd %+v", err)
	}
	temp.Close()
	t, tErr := os.Open("temp")
	if tErr != nil {
		return 0, 0, fmt.Errorf("error during Opening the list file %+v", tErr)
	}
	r := bufio.NewScanner(t)
	r.Split(bufio.ScanLines)

	var (
		f         *os.File
		err       error
		filesRead fileRead
		bytesRead byteRead
	)
	b := bytes.NewBuffer(nil)
	for r.Scan() {
		if (len(sigCh))==1{break}
		f, err = os.Open("./json/" + r.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}
		io.Copy(b, f)
		f.Close()
		jsondelete.DeleteFile("./json/"+r.Text())
		filesRead += 1
		bytesRead += byteRead(b.Len())
		js := &jsonwriter.Js{}
		json.Unmarshal(b.Bytes(),js)
	}
	t.Close()
	jsondelete.DeleteFile("temp")
	return filesRead, bytesRead, err
}
