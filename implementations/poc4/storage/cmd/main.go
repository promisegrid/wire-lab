package main

import (
	"context"
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc4/storage"
)

func main() {
	if err := storage.Main(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
