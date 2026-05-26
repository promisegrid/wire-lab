package main

import (
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc5/relay"
)

func main() {
	if err := relay.Command(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
