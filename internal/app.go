package internal

import (
	"github.com/praveenmahasena/generate_json/internal/cli"
	"github.com/praveenmahasena/generate_json/internal/jsonwriter"
)

func Run() error {
	amount := cli.ReadAmount()
	js := jsonwriter.New(amount)
	return js.Generate()
}
