package main

import (
	"fmt"
	"os"

	"github.com/praveenmahasena/generate_json/internal"
)

func main(){
	if err:=internal.GracefulRead();err!=nil{
		fmt.Fprintln(os.Stderr,err)
	}
}
