package main

import (
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc4/relay"
)

func main() {
	if err := relay.Command(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
