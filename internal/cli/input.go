package cli

import (
	"flag"
)

var (
	defaultAmount uint=10_000
)

func ReadAmount() uint {
	amount := flag.Uint("amount", defaultAmount, "amount of json files this binary should generate would be defaulted to 10_000")
	flag.Parse()
	return *(amount) // call me old idc but been doing Clang for awhile and this *() reference seems more appealing then just doing *
}
