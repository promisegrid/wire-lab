package main

import (
	"encoding/json"
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/analyzer"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: poc17-analyze <run-dir>")
		os.Exit(2)
	}
	summary, err := analyzer.Analyze(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "analyze: %v\n", err)
		os.Exit(1)
	}
	encoded, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode summary: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(encoded))
}
