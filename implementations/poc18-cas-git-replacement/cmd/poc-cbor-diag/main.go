// Command poc-cbor-diag renders exact POC18 CBOR messages for review.
//
// Intent: Raw message diagnostics are a development aid for checking pCID,
// parent, payload, and proof shape. Source: DI-harih
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/scenario"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

type fixtureDiagnosticResult struct {
	workspace.IngestResult
	AliceStoreRoot string          `json:"alice_store_root"`
	Scenario       scenario.Result `json:"scenario"`
}

type messageDAGRecord struct {
	Source      string `json:"source"`
	Observer    string `json:"observer"`
	Direction   string `json:"direction"`
	Peer        string `json:"peer"`
	Protocol    string `json:"protocol"`
	PromiseKind string `json:"promise_kind"`
	ExactCID    string `json:"exact_cid"`
	Path        string `json:"path"`
}

type diagnosticIndexEntry struct {
	Flow           string `json:"flow"`
	CID            string `json:"cid"`
	Source         string `json:"source"`
	PromiseKind    string `json:"promise_kind"`
	ArtifactPath   string `json:"artifact_path"`
	DiagnosticPath string `json:"diagnostic_path"`
}

type diagnosticReportBuilder struct {
	runRoot     string
	observerRun string
	outDir      string
	fixture     fixtureDiagnosticResult
	aliceCAS    *store.FileStore
	index       []diagnosticIndexEntry
}

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
	diagnosticReport := flag.Bool("diagnostic-report", false, "write named diagnostic files for the deterministic POC18 run")
	outDir := flag.String("out-dir", "", "diagnostic report output directory")
	flag.Parse()
	if *diagnosticReport {
		return writeDiagnosticReport(*runRoot, *outDir)
	}
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

func writeDiagnosticReport(runRoot string, outDir string) error {
	if runRoot == "" {
		return fmt.Errorf("-run-root is required with -diagnostic-report")
	}
	if outDir == "" {
		outDir = filepath.Join(runRoot, "poc18-diagnostics")
	}
	fixture, fixtureErr := readFixtureDiagnosticResult(filepath.Join(runRoot, "result.json"))
	if fixtureErr != nil {
		return fixtureErr
	}
	aliceCAS, openErr := store.Open(fixture.AliceStoreRoot)
	if openErr != nil {
		return openErr
	}
	builder := &diagnosticReportBuilder{
		runRoot:     runRoot,
		observerRun: filepath.Join(runRoot, "observer", "poc18-demo"),
		outDir:      outDir,
		fixture:     fixture,
		aliceCAS:    aliceCAS,
	}
	return builder.write()
}

func readFixtureDiagnosticResult(path string) (fixtureDiagnosticResult, error) {
	var result fixtureDiagnosticResult
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return fixtureDiagnosticResult{}, readErr
	}
	if err := json.Unmarshal(content, &result); err != nil {
		return fixtureDiagnosticResult{}, err
	}
	if result.AliceStoreRoot == "" {
		return fixtureDiagnosticResult{}, fmt.Errorf("result.json is missing alice_store_root")
	}
	return result, nil
}

func (builder *diagnosticReportBuilder) write() error {
	if mkdirErr := os.MkdirAll(builder.outDir, 0o755); mkdirErr != nil {
		return mkdirErr
	}
	// Intent: The report names fixed review flows so `run-clean.sh` can fail when
	// POC18 stops producing one of the core version-control or peer-fetch message
	// shapes. Source: DI-basan
	storeSelections := []struct {
		flow        string
		cidText     string
		source      string
		promiseKind string
	}{
		{flow: "reference-set", cidText: builder.fixture.Scenario.MergeBranchRefSetCID, source: "result.scenario.merge_branch_reference_set_cid", promiseKind: "reference_set"},
		{flow: "node-version", cidText: builder.fixture.Scenario.RenameLabelNodeCID, source: "result.scenario.rename_label_node_cid", promiseKind: "posix_node"},
		{flow: "snapshot", cidText: builder.fixture.Scenario.RenameCopySnapshotCID, source: "result.scenario.rename_copy_snapshot_cid", promiseKind: "snapshot"},
		{flow: "review-statement", cidText: builder.fixture.Scenario.TestStatementCID, source: "result.scenario.test_statement_cid", promiseKind: "review_statement"},
		{flow: "merge-snapshot", cidText: builder.fixture.Scenario.MergeSnapshotCID, source: "result.scenario.merge_snapshot_cid", promiseKind: "snapshot"},
	}
	for _, selection := range storeSelections {
		if err := builder.renderStoreCID(selection.flow, selection.cidText, selection.source, selection.promiseKind); err != nil {
			return err
		}
	}
	if err := builder.renderStorePredicate("directory-node", "alice sparse CAS posix_node directory scan", "posix_node", isDirectoryNode); err != nil {
		return err
	}
	if err := builder.renderStorePredicate("materialization-object", "alice sparse CAS chunk_manifest scan", "chunk_manifest", anyEnvelope); err != nil {
		return err
	}
	for _, selection := range []struct {
		flow        string
		promiseKind string
	}{
		{flow: "sync-interest", promiseKind: "sync_interest"},
		{flow: "object-availability", promiseKind: "object_availability"},
		{flow: "object-retrieval-redemption", promiseKind: "object_retrieval_redemption"},
	} {
		if err := builder.renderCollectorMessage(selection.flow, selection.promiseKind); err != nil {
			return err
		}
	}
	return builder.writeIndex()
}

func (builder *diagnosticReportBuilder) renderStoreCID(flow string, cidText string, source string, promiseKind string) error {
	if cidText == "" {
		return fmt.Errorf("%s CID is missing", flow)
	}
	objectCID, parseErr := store.ParseCIDText(cidText)
	if parseErr != nil {
		return parseErr
	}
	content, entry, getErr := builder.aliceCAS.Get(objectCID)
	if getErr != nil {
		return getErr
	}
	return builder.renderBytes(flow, cidText, source, promiseKind, entry.Path, content)
}

func (builder *diagnosticReportBuilder) renderStorePredicate(flow string, source string, promiseKind string, predicate func(graph.EnvelopeView) bool) error {
	entries, listErr := builder.aliceCAS.List()
	if listErr != nil {
		return listErr
	}
	for _, entry := range entries {
		if entry.Kind != "message" {
			continue
		}
		objectCID, parseErr := store.ParseCIDText(entry.CID)
		if parseErr != nil {
			return parseErr
		}
		content, _, getErr := builder.aliceCAS.Get(objectCID)
		if getErr != nil {
			return getErr
		}
		view, parseEnvelopeErr := graph.ParseEnvelope(content)
		if parseEnvelopeErr != nil {
			continue
		}
		kind, kindErr := view.PayloadKind()
		if kindErr != nil || kind != promiseKind || !predicate(view) {
			continue
		}
		return builder.renderBytes(flow, entry.CID, source, promiseKind, entry.Path, content)
	}
	return fmt.Errorf("no %s diagnostic candidate found", flow)
}

func isDirectoryNode(view graph.EnvelopeView) bool {
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil || len(body) < 2 {
		return false
	}
	nodeType, ok := body[1].(string)
	return ok && nodeType == "directory"
}

func anyEnvelope(graph.EnvelopeView) bool {
	return true
}

func (builder *diagnosticReportBuilder) renderCollectorMessage(flow string, promiseKind string) error {
	record, recordErr := builder.firstCollectorMessage(promiseKind)
	if recordErr != nil {
		return recordErr
	}
	artifactPath := filepath.Join(builder.observerRun, filepath.FromSlash(record.Path))
	content, readErr := os.ReadFile(artifactPath)
	if readErr != nil {
		return readErr
	}
	return builder.renderBytes(flow, record.ExactCID, "observer message-dag "+record.Direction+" "+record.Observer+"->"+record.Peer, promiseKind, artifactPath, content)
}

func (builder *diagnosticReportBuilder) firstCollectorMessage(promiseKind string) (messageDAGRecord, error) {
	indexPath := filepath.Join(builder.observerRun, "message-dag.jsonl")
	file, openErr := os.Open(indexPath)
	if openErr != nil {
		return messageDAGRecord{}, openErr
	}
	defer closeFile(file, "message-dag")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record messageDAGRecord
		if unmarshalErr := json.Unmarshal(scanner.Bytes(), &record); unmarshalErr != nil {
			return messageDAGRecord{}, unmarshalErr
		}
		if record.PromiseKind == promiseKind {
			return record, nil
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return messageDAGRecord{}, scanErr
	}
	return messageDAGRecord{}, fmt.Errorf("message-dag missing promise kind %s", promiseKind)
}

func (builder *diagnosticReportBuilder) renderBytes(flow string, cidText string, source string, promiseKind string, artifactPath string, content []byte) error {
	diagnostic, diagErr := graph.Diagnostic(content)
	if diagErr != nil {
		return diagErr
	}
	diagnosticPath := filepath.Join(builder.outDir, flow+".diag.txt")
	rendered := fmt.Sprintf("flow: %s\ncid: %s\nsource: %s\npromise_kind: %s\nartifact_path: %s\n\n%s", flow, cidText, source, promiseKind, artifactPath, diagnostic)
	if writeErr := os.WriteFile(diagnosticPath, []byte(rendered), 0o644); writeErr != nil {
		return writeErr
	}
	builder.index = append(builder.index, diagnosticIndexEntry{
		Flow:           flow,
		CID:            cidText,
		Source:         source,
		PromiseKind:    promiseKind,
		ArtifactPath:   artifactPath,
		DiagnosticPath: diagnosticPath,
	})
	return nil
}

func (builder *diagnosticReportBuilder) writeIndex() error {
	indexPath := filepath.Join(builder.outDir, "index.json")
	content, marshalErr := json.MarshalIndent(builder.index, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	content = append(content, '\n')
	return os.WriteFile(indexPath, content, 0o644)
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

func closeFile(file *os.File, label string) {
	if closeErr := file.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "close %s file: %v\n", label, closeErr)
	}
}
