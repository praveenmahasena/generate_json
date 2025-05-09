package cli

import (
	"flag"
)

func ReadAmount() uint {
	amount := flag.Uint("amount", 10_000, "amount of json files this binary should generate would be defaulted to 10_000")
	flag.Parse()
	return *(amount) // been doing Clang for awhile and this *() reference seems correct then just doing *
}
