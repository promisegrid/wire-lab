package main

import (
	"context"
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc4/signed"
)

func main() {
	if err := signed.Main(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
