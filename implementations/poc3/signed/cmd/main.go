package main

import (
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc3/signed"
)

func main() {
	if err := signed.Command(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
