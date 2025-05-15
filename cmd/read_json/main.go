package main

import (
	"fmt"
	"os"

	"github.com/praveenmahasena/generate_json/internal/jsonreader"
	"github.com/praveenmahasena/generate_json/internal/logger"
)

func main() {
	fileRead, byteRead, err := jsonreader.Read(logger.Logger())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// there is no a return statement here it's not an error
		// I just wanna show it as the file read are 0
	}
	fmt.Printf("total files read %v \n", fileRead)
	fmt.Printf("total bytes read %v \n", byteRead)
}
