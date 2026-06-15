package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const defaultRunRoot = "/run/poc15/poc15-demo"

// diagnosticValue is the small CBOR value tree used by the POC15 diagnostic
// tool.
// Intent: Operators need a read-only way to inspect exact message artifacts
// without converting those artifacts into protocol behavior or agent-visible
// coordination. Source: DI-bapif
type diagnosticValue struct {
	kind     string
	unsigned uint64
	negative int64
	bytes    []byte
	text     string
	array    []diagnosticValue
	mapPairs []diagnosticPair
	tag      uint64
	tagValue *diagnosticValue
	simple   byte
}

type diagnosticPair struct {
	key   diagnosticValue
	value diagnosticValue
}

type diagnosticOptions struct {
	expandBytes bool
	maxBytes    int
}

// diagnosticDecoder parses and renders the deterministic CBOR subset used by
// POC15, plus enough generic CBOR to make malformed or future artifacts clear.
// Intent: Keep the inspector self-contained in Go and free of third-party
// dependencies while still producing diagnostic notation for retained raw
// message bytes. Source: DI-bapif
type diagnosticDecoder struct {
	data   []byte
	offset int
}

func newDiagnosticDecoder(data []byte) *diagnosticDecoder {
	return &diagnosticDecoder{data: data}
}

func (decoder *diagnosticDecoder) parseOne() (diagnosticValue, error) {
	initialByte, byteErr := decoder.readByte()
	if byteErr != nil {
		return diagnosticValue{}, byteErr
	}
	majorType := initialByte >> 5
	additionalInfo := initialByte & 0x1f
	length, lengthErr := decoder.readLength(additionalInfo)
	if lengthErr != nil {
		return diagnosticValue{}, lengthErr
	}
	switch majorType {
	case 0:
		return diagnosticValue{kind: "uint", unsigned: length.value}, nil
	case 1:
		if length.value > uint64(^uint64(0)>>1) {
			return diagnosticValue{}, fmt.Errorf("negative integer is too large: %d", length.value)
		}
		return diagnosticValue{kind: "nint", negative: -1 - int64(length.value)}, nil
	case 2:
		if length.indefinite {
			return decoder.parseIndefiniteBytes()
		}
		byteString, readErr := decoder.readBytes(int(length.value))
		if readErr != nil {
			return diagnosticValue{}, readErr
		}
		return diagnosticValue{kind: "bytes", bytes: byteString}, nil
	case 3:
		if length.indefinite {
			return decoder.parseIndefiniteText()
		}
		textBytes, readErr := decoder.readBytes(int(length.value))
		if readErr != nil {
			return diagnosticValue{}, readErr
		}
		if !utf8.Valid(textBytes) {
			return diagnosticValue{}, fmt.Errorf("invalid utf-8 text string")
		}
		return diagnosticValue{kind: "text", text: string(textBytes)}, nil
	case 4:
		if length.indefinite {
			return decoder.parseIndefiniteArray()
		}
		items := make([]diagnosticValue, 0, int(length.value))
		for itemIndex := uint64(0); itemIndex < length.value; itemIndex++ {
			item, itemErr := decoder.parseOne()
			if itemErr != nil {
				return diagnosticValue{}, itemErr
			}
			items = append(items, item)
		}
		return diagnosticValue{kind: "array", array: items}, nil
	case 5:
		if length.indefinite {
			return decoder.parseIndefiniteMap()
		}
		pairs := make([]diagnosticPair, 0, int(length.value))
		for pairIndex := uint64(0); pairIndex < length.value; pairIndex++ {
			key, keyErr := decoder.parseOne()
			if keyErr != nil {
				return diagnosticValue{}, keyErr
			}
			value, valueErr := decoder.parseOne()
			if valueErr != nil {
				return diagnosticValue{}, valueErr
			}
			pairs = append(pairs, diagnosticPair{key: key, value: value})
		}
		return diagnosticValue{kind: "map", mapPairs: pairs}, nil
	case 6:
		if length.indefinite {
			return diagnosticValue{}, fmt.Errorf("tag cannot have indefinite length")
		}
		tagValue, tagErr := decoder.parseOne()
		if tagErr != nil {
			return diagnosticValue{}, tagErr
		}
		return diagnosticValue{kind: "tag", tag: length.value, tagValue: &tagValue}, nil
	case 7:
		return decoder.parseSimple(additionalInfo, length)
	default:
		return diagnosticValue{}, fmt.Errorf("unsupported CBOR major type %d", majorType)
	}
}

func (decoder *diagnosticDecoder) parseAll() (diagnosticValue, error) {
	value, valueErr := decoder.parseOne()
	if valueErr != nil {
		return diagnosticValue{}, valueErr
	}
	if decoder.offset != len(decoder.data) {
		return diagnosticValue{}, fmt.Errorf("trailing CBOR bytes: %d", len(decoder.data)-decoder.offset)
	}
	return value, nil
}

func (decoder *diagnosticDecoder) parseIndefiniteBytes() (diagnosticValue, error) {
	var chunks []byte
	for !decoder.consumeBreak() {
		item, itemErr := decoder.parseOne()
		if itemErr != nil {
			return diagnosticValue{}, itemErr
		}
		if item.kind != "bytes" {
			return diagnosticValue{}, fmt.Errorf("indefinite byte string chunk has kind %s", item.kind)
		}
		chunks = append(chunks, item.bytes...)
	}
	return diagnosticValue{kind: "bytes", bytes: chunks}, nil
}

func (decoder *diagnosticDecoder) parseIndefiniteText() (diagnosticValue, error) {
	var builder strings.Builder
	for !decoder.consumeBreak() {
		item, itemErr := decoder.parseOne()
		if itemErr != nil {
			return diagnosticValue{}, itemErr
		}
		if item.kind != "text" {
			return diagnosticValue{}, fmt.Errorf("indefinite text string chunk has kind %s", item.kind)
		}
		builder.WriteString(item.text)
	}
	return diagnosticValue{kind: "text", text: builder.String()}, nil
}

func (decoder *diagnosticDecoder) parseIndefiniteArray() (diagnosticValue, error) {
	var items []diagnosticValue
	for !decoder.consumeBreak() {
		item, itemErr := decoder.parseOne()
		if itemErr != nil {
			return diagnosticValue{}, itemErr
		}
		items = append(items, item)
	}
	return diagnosticValue{kind: "array", array: items}, nil
}

func (decoder *diagnosticDecoder) parseIndefiniteMap() (diagnosticValue, error) {
	var pairs []diagnosticPair
	for !decoder.consumeBreak() {
		key, keyErr := decoder.parseOne()
		if keyErr != nil {
			return diagnosticValue{}, keyErr
		}
		value, valueErr := decoder.parseOne()
		if valueErr != nil {
			return diagnosticValue{}, valueErr
		}
		pairs = append(pairs, diagnosticPair{key: key, value: value})
	}
	return diagnosticValue{kind: "map", mapPairs: pairs}, nil
}

func (decoder *diagnosticDecoder) parseSimple(additionalInfo byte, length cborLength) (diagnosticValue, error) {
	switch additionalInfo {
	case 20:
		return diagnosticValue{kind: "simple", simple: 20}, nil
	case 21:
		return diagnosticValue{kind: "simple", simple: 21}, nil
	case 22:
		return diagnosticValue{kind: "simple", simple: 22}, nil
	case 23:
		return diagnosticValue{kind: "simple", simple: 23}, nil
	case 24:
		if length.indefinite || length.value > 255 {
			return diagnosticValue{}, fmt.Errorf("invalid simple value")
		}
		return diagnosticValue{kind: "simple", simple: byte(length.value)}, nil
	default:
		return diagnosticValue{kind: "simple", simple: additionalInfo}, nil
	}
}

type cborLength struct {
	value      uint64
	indefinite bool
}

func (decoder *diagnosticDecoder) readLength(additionalInfo byte) (cborLength, error) {
	switch additionalInfo {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23:
		return cborLength{value: uint64(additionalInfo)}, nil
	case 24:
		value, readErr := decoder.readByte()
		if readErr != nil {
			return cborLength{}, readErr
		}
		return cborLength{value: uint64(value)}, nil
	case 25:
		value, readErr := decoder.readUint(2)
		return cborLength{value: value}, readErr
	case 26:
		value, readErr := decoder.readUint(4)
		return cborLength{value: value}, readErr
	case 27:
		value, readErr := decoder.readUint(8)
		return cborLength{value: value}, readErr
	case 31:
		return cborLength{indefinite: true}, nil
	default:
		return cborLength{}, fmt.Errorf("unsupported additional information %d", additionalInfo)
	}
}

func (decoder *diagnosticDecoder) readByte() (byte, error) {
	if decoder.offset >= len(decoder.data) {
		return 0, io.ErrUnexpectedEOF
	}
	value := decoder.data[decoder.offset]
	decoder.offset++
	return value, nil
}

func (decoder *diagnosticDecoder) readBytes(length int) ([]byte, error) {
	if length < 0 || decoder.offset+length > len(decoder.data) {
		return nil, io.ErrUnexpectedEOF
	}
	value := make([]byte, length)
	copy(value, decoder.data[decoder.offset:decoder.offset+length])
	decoder.offset += length
	return value, nil
}

func (decoder *diagnosticDecoder) readUint(length int) (uint64, error) {
	bytes, readErr := decoder.readBytes(length)
	if readErr != nil {
		return 0, readErr
	}
	var value uint64
	for _, byteValue := range bytes {
		value = value<<8 | uint64(byteValue)
	}
	return value, nil
}

func (decoder *diagnosticDecoder) consumeBreak() bool {
	if decoder.offset < len(decoder.data) && decoder.data[decoder.offset] == 0xff {
		decoder.offset++
		return true
	}
	return false
}

func (value diagnosticValue) render(options diagnosticOptions) string {
	return value.renderAt(0, options)
}

func (value diagnosticValue) renderAt(indent int, options diagnosticOptions) string {
	switch value.kind {
	case "uint":
		return fmt.Sprintf("%d", value.unsigned)
	case "nint":
		return fmt.Sprintf("%d", value.negative)
	case "bytes":
		return value.renderBytes(indent, options)
	case "text":
		return quoteDiagnosticText(value.text)
	case "array":
		return renderList("[", "]", value.renderedArrayItems(indent, options), indent)
	case "map":
		return renderList("{", "}", value.renderedMapItems(indent, options), indent)
	case "tag":
		return fmt.Sprintf("%d(%s)", value.tag, value.tagValue.renderAt(indent, options))
	case "simple":
		return renderSimple(value.simple)
	default:
		return fmt.Sprintf("<unknown:%s>", value.kind)
	}
}

func (value diagnosticValue) renderBytes(indent int, options diagnosticOptions) string {
	if options.expandBytes {
		nestedValue, nestedErr := parseDiagnosticValue(value.bytes)
		if nestedErr == nil {
			return "<< " + nestedValue.renderAt(indent+1, options) + " >>"
		}
	}
	hexBytes := hex.EncodeToString(value.bytes)
	if options.maxBytes > 0 && len(value.bytes) > options.maxBytes {
		hexBytes = hexBytes[:options.maxBytes*2] + "..."
	}
	return "h'" + hexBytes + "'"
}

func (value diagnosticValue) renderedArrayItems(indent int, options diagnosticOptions) []string {
	renderedItems := make([]string, 0, len(value.array))
	for _, item := range value.array {
		renderedItems = append(renderedItems, item.renderAt(indent+1, options))
	}
	return renderedItems
}

func (value diagnosticValue) renderedMapItems(indent int, options diagnosticOptions) []string {
	renderedItems := make([]string, 0, len(value.mapPairs))
	for _, pair := range value.mapPairs {
		renderedItems = append(renderedItems, pair.key.renderAt(indent+1, options)+": "+pair.value.renderAt(indent+1, options))
	}
	return renderedItems
}

func renderList(openDelimiter, closeDelimiter string, renderedItems []string, indent int) string {
	if len(renderedItems) == 0 {
		return openDelimiter + closeDelimiter
	}
	var builder strings.Builder
	builder.WriteString(openDelimiter)
	builder.WriteString("\n")
	for _, renderedItem := range renderedItems {
		builder.WriteString(strings.Repeat("  ", indent+1))
		builder.WriteString(renderedItem)
		builder.WriteString(",\n")
	}
	builder.WriteString(strings.Repeat("  ", indent))
	builder.WriteString(closeDelimiter)
	return builder.String()
}

func renderSimple(simple byte) string {
	switch simple {
	case 20:
		return "false"
	case 21:
		return "true"
	case 22:
		return "null"
	case 23:
		return "undefined"
	default:
		return fmt.Sprintf("simple(%d)", simple)
	}
}

func quoteDiagnosticText(text string) string {
	encoded, marshalErr := json.Marshal(text)
	if marshalErr != nil {
		return fmt.Sprintf("%q", text)
	}
	return string(encoded)
}

func parseDiagnosticValue(data []byte) (diagnosticValue, error) {
	decoder := newDiagnosticDecoder(data)
	return decoder.parseAll()
}

type messageDAGRecord struct {
	Source            string `json:"source"`
	Observer          string `json:"observer"`
	Direction         string `json:"direction"`
	Peer              string `json:"peer"`
	Protocol          string `json:"protocol"`
	ExactSHA256       string `json:"exact_sha256"`
	ParentExactSHA256 string `json:"parent_exact_sha256,omitempty"`
	PromiseAbout      string `json:"promise_about"`
	SourceEvent       string `json:"source_event"`
	Path              string `json:"path"`
}

type hashList []string

func (hashes *hashList) String() string {
	return strings.Join(*hashes, ",")
}

func (hashes *hashList) Set(value string) error {
	if value == "" {
		return fmt.Errorf("hash cannot be empty")
	}
	*hashes = append(*hashes, value)
	return nil
}

type cliConfig struct {
	runRoot     string
	hashes      hashList
	expandBytes bool
	maxBytes    int
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "poc15-cbor-diag: %v\n", err)
		os.Exit(1)
	}
}

func run(arguments []string, output io.Writer) error {
	config, paths, parseErr := parseCLI(arguments)
	if parseErr != nil {
		return parseErr
	}
	selectedPaths, selectErr := selectInputPaths(config, paths)
	if selectErr != nil {
		return selectErr
	}
	if len(selectedPaths) == 0 {
		return errors.New("provide one or more paths or -hash values")
	}
	options := diagnosticOptions{expandBytes: config.expandBytes, maxBytes: config.maxBytes}
	for pathIndex, path := range selectedPaths {
		if pathIndex > 0 {
			if writeErr := writeLine(output, ""); writeErr != nil {
				return writeErr
			}
		}
		if err := renderPath(path, config.runRoot, options, output); err != nil {
			return err
		}
	}
	return nil
}

func parseCLI(arguments []string) (cliConfig, []string, error) {
	config := cliConfig{runRoot: defaultRunRoot, expandBytes: true, maxBytes: 96}
	flagSet := flag.NewFlagSet("poc15-cbor-diag", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	flagSet.StringVar(&config.runRoot, "run-root", config.runRoot, "POC15 run root containing message-dag.jsonl and message-cas/")
	flagSet.Var(&config.hashes, "hash", "exact_sha256 of a retained message artifact; may be repeated")
	flagSet.BoolVar(&config.expandBytes, "expand-bytes", config.expandBytes, "render nested CBOR byte strings as diagnostic notation")
	flagSet.IntVar(&config.maxBytes, "max-bytes", config.maxBytes, "maximum raw byte-string length to show before truncation; zero disables truncation")
	if err := flagSet.Parse(arguments); err != nil {
		return cliConfig{}, nil, err
	}
	return config, flagSet.Args(), nil
}

func selectInputPaths(config cliConfig, directPaths []string) ([]string, error) {
	var selectedPaths []string
	for _, hashValue := range config.hashes {
		artifactPath, pathErr := pathForHash(config.runRoot, hashValue)
		if pathErr != nil {
			return nil, pathErr
		}
		selectedPaths = append(selectedPaths, artifactPath)
	}
	selectedPaths = append(selectedPaths, directPaths...)
	return selectedPaths, nil
}

func pathForHash(runRoot, hashValue string) (string, error) {
	record, recordErr := findMessageDAGRecord(runRoot, hashValue)
	if recordErr != nil {
		return "", recordErr
	}
	return filepath.Join(runRoot, record.Path), nil
}

func findMessageDAGRecord(runRoot, hashValue string) (messageDAGRecord, error) {
	indexPath := filepath.Join(runRoot, "message-dag.jsonl")
	indexFile, openErr := os.Open(indexPath)
	if openErr != nil {
		return messageDAGRecord{}, openErr
	}
	defer func() {
		closeErr := indexFile.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc15-cbor-diag: close %s: %v\n", indexPath, closeErr)
		}
	}()
	scanner := bufio.NewScanner(indexFile)
	for scanner.Scan() {
		var record messageDAGRecord
		if unmarshalErr := json.Unmarshal(scanner.Bytes(), &record); unmarshalErr != nil {
			return messageDAGRecord{}, unmarshalErr
		}
		if record.ExactSHA256 == hashValue {
			return record, nil
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return messageDAGRecord{}, scanErr
	}
	return messageDAGRecord{}, fmt.Errorf("hash %s not found in %s", hashValue, indexPath)
}

func renderPath(path, runRoot string, options diagnosticOptions, output io.Writer) error {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return readErr
	}
	value, parseErr := parseDiagnosticValue(data)
	if parseErr != nil {
		return parseErr
	}
	metadataLine, metadataErr := metadataForPath(runRoot, path)
	if metadataErr != nil {
		if writeErr := writeFormatted(output, "# metadata: unavailable: %v\n", metadataErr); writeErr != nil {
			return writeErr
		}
	} else if metadataLine != "" {
		if writeErr := writeFormatted(output, "# metadata: %s\n", metadataLine); writeErr != nil {
			return writeErr
		}
	}
	if writeErr := writeFormatted(output, "# path: %s\n", path); writeErr != nil {
		return writeErr
	}
	if writeErr := writeFormatted(output, "# bytes: %d\n", len(data)); writeErr != nil {
		return writeErr
	}
	if writeErr := writeLine(output, value.render(options)); writeErr != nil {
		return writeErr
	}
	return nil
}

func writeFormatted(output io.Writer, format string, arguments ...any) error {
	_, writeErr := fmt.Fprintf(output, format, arguments...)
	return writeErr
}

func writeLine(output io.Writer, line string) error {
	_, writeErr := fmt.Fprintln(output, line)
	return writeErr
}

func metadataForPath(runRoot, path string) (string, error) {
	relativePath, relativeErr := filepath.Rel(runRoot, path)
	if relativeErr != nil || strings.HasPrefix(relativePath, "..") {
		return "", relativeErr
	}
	indexPath := filepath.Join(runRoot, "message-dag.jsonl")
	indexFile, openErr := os.Open(indexPath)
	if openErr != nil {
		return "", openErr
	}
	defer func() {
		closeErr := indexFile.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc15-cbor-diag: close %s: %v\n", indexPath, closeErr)
		}
	}()
	scanner := bufio.NewScanner(indexFile)
	for scanner.Scan() {
		var record messageDAGRecord
		if unmarshalErr := json.Unmarshal(scanner.Bytes(), &record); unmarshalErr != nil {
			return "", unmarshalErr
		}
		if record.Path == relativePath {
			return formatMessageDAGRecord(record), nil
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return "", scanErr
	}
	return "", nil
}

func formatMessageDAGRecord(record messageDAGRecord) string {
	fields := []string{
		"source=" + record.Source,
		"observer=" + record.Observer,
		"direction=" + record.Direction,
		"peer=" + record.Peer,
		"protocol=" + record.Protocol,
		"exact_sha256=" + record.ExactSHA256,
		"promise_about=" + record.PromiseAbout,
	}
	if record.ParentExactSHA256 != "" {
		fields = append(fields, "parent_exact_sha256="+record.ParentExactSHA256)
	}
	return strings.Join(fields, " ")
}
