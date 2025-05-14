package main

import (
	"fmt"
	"os"

	"github.com/praveenmahasena/generate_json/internal/cli"
	"github.com/praveenmahasena/generate_json/internal/jsonwriter"
	"github.com/praveenmahasena/generate_json/internal/logger"
	//"syscall"
)

func main() {
	amount := cli.ReadAmount()
	js := jsonwriter.New(amount,logger.Logger())
	if err:= js.Generate();err!=nil{
		fmt.Fprintln(os.Stderr,err)
	}
}
