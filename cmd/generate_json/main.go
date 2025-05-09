package main

import (
	"fmt"
	"os"
	//"syscall"

	"github.com/praveenmahasena/generate_json/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		// syscall.Exit(1) // I do know I do not need this line here but I'm just so used to having it here
	}
}
