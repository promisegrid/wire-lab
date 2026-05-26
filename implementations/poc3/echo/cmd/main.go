package main

import (
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc3/echo"
)

func main() {
	if err := echo.Command(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
