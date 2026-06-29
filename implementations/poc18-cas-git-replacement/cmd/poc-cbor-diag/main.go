// Command poc-cbor-diag renders exact POC18 CBOR messages for review.
//
// Intent: Raw message diagnostics are a development aid for checking pCID,
// parent, payload, and proof shape. Source: DI-harih
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "poc-cbor-diag: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	runRoot := flag.String("run-root", "", "read result.json and render its diagnostic message")
	storeRoot := flag.String("store", "", "CAS store root")
	cidText := flag.String("cid", "", "CID to render from CAS")
	filePath := flag.String("file", "", "CBOR file to render")
	flag.Parse()
	if *runRoot != "" {
		result, resultErr := readResult(filepath.Join(*runRoot, "result.json"))
		if resultErr != nil {
			return resultErr
		}
		*storeRoot = result.StoreRoot
		*cidText = result.DiagnosticMessageCID
	}
	content, loadErr := loadBytes(*storeRoot, *cidText, *filePath)
	if loadErr != nil {
		return loadErr
	}
	diagnostic, diagErr := graph.Diagnostic(content)
	if diagErr != nil {
		return diagErr
	}
	fmt.Print(diagnostic)
	return nil
}

func readResult(path string) (workspace.IngestResult, error) {
	var result workspace.IngestResult
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return workspace.IngestResult{}, readErr
	}
	if err := json.Unmarshal(content, &result); err != nil {
		return workspace.IngestResult{}, err
	}
	return result, nil
}

func loadBytes(storeRoot, cidText, filePath string) ([]byte, error) {
	if filePath != "" {
		return os.ReadFile(filePath)
	}
	if storeRoot == "" || cidText == "" {
		return nil, fmt.Errorf("either -file, -run-root, or both -store and -cid are required")
	}
	cas, openErr := store.Open(storeRoot)
	if openErr != nil {
		return nil, openErr
	}
	objectCID, parseErr := store.ParseCIDText(cidText)
	if parseErr != nil {
		return nil, parseErr
	}
	content, _, getErr := cas.Get(objectCID)
	return content, getErr
}
