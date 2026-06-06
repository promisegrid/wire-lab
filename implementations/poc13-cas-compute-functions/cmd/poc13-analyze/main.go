package main

import (
	"encoding/json"
	"fmt"
	"os"

	poc13 "promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: poc13-analyze RUN_DIR\n")
		os.Exit(2)
	}
	summary, err := poc13.AnalyzeRun(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "poc13-analyze: %v\n", err)
		os.Exit(1)
	}
	if err := poc13.ValidateAnalysis(summary); err != nil {
		fmt.Fprintf(os.Stderr, "poc13-analyze: acceptance criteria failed: %v\n", err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		fmt.Fprintf(os.Stderr, "poc13-analyze: encode summary: %v\n", err)
		os.Exit(1)
	}
}
