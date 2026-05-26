package main

import (
	"context"
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc5/fibonacci"
)

func main() {
	if err := fibonacci.Main(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
