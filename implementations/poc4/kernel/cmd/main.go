package main

import (
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc4/kernel"
)

func main() {
	if err := kernel.Command(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
