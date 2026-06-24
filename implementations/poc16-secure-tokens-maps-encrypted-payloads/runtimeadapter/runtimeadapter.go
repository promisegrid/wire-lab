package runtimeadapter

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/tetratelabs/wazero"
)

const (
	PromiseAboutWASMAdapter       = "wasm_adapter_event"
	PromiseAboutStdioAdapter      = "stdio_adapter_event"
	PromiseAboutWASMModuleUse     = "wasm_module_execution"
	PromiseAboutStdioWorkerUse    = "stdio_worker_roundtrip"
	PromiseAboutLocalEventSummary = "local_event_summary"
	PromiseAboutPeerAttestation   = "peer_carried_attestation"
	PromiseAboutExchangeRate      = "bearer_token_exchange_rate"
	PromiseAboutTopologySignal    = "relationship_topology_signal"
	PromiseAboutVoluntaryGossip   = "voluntary_gossip"

	WASMExportName     = "promise_fibonacci"
	ExpectedWASMInput  = 9
	ExpectedWASMResult = 34

	maxStdioCBORFrameSize = 16 * 1024 * 1024
)

// MinimalWASMModule is a deterministic no-import WebAssembly module that
// exports promise_fibonacci(i32) -> i32.
// Intent: Peggy must execute real WASM with wazero, not merely validate module
// header bytes, while still keeping PromiseGrid semantics in ordinary
// pCID-defined envelopes outside the sandbox. Source: DI-kimim; DI-sivis
var MinimalWASMModule = []byte{
	0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00,
	0x01, 0x06, 0x01, 0x60, 0x01, 0x7f, 0x01, 0x7f,
	0x03, 0x02, 0x01, 0x00,
	0x07, 0x15, 0x01, 0x11, 0x70, 0x72, 0x6f, 0x6d, 0x69, 0x73, 0x65, 0x5f, 0x66, 0x69, 0x62, 0x6f, 0x6e, 0x61, 0x63, 0x63, 0x69, 0x00, 0x00,
	0x0a, 0x1e, 0x01, 0x1c, 0x00,
	0x20, 0x00, 0x41, 0x02, 0x48, 0x04, 0x7f, 0x20, 0x00, 0x05,
	0x20, 0x00, 0x41, 0x01, 0x6b, 0x10, 0x00, 0x20, 0x00, 0x41, 0x02, 0x6b, 0x10, 0x00, 0x6a, 0x0b, 0x0b,
}

// WASMRunResult records the small amount of deterministic runtime output that
// Peggy can promise about her local WASM execution without turning wazero into a
// PromiseGrid command surface.
type WASMRunResult struct {
	ExportName  string
	InputValue  uint64
	ExportValue uint64
}

// ValidateWASMModule checks the stable magic and version bytes before the real
// wazero compile/instantiate/call path runs. This is only a preflight error
// message helper; RunWASMModule is the behavior that proves execution.
func ValidateWASMModule(moduleBytes []byte) error {
	if len(moduleBytes) < 8 {
		return fmt.Errorf("wasm module too short: %d bytes", len(moduleBytes))
	}
	if !bytes.Equal(moduleBytes[:4], MinimalWASMModule[:4]) {
		return fmt.Errorf("wasm module magic mismatch")
	}
	if !bytes.Equal(moduleBytes[4:8], MinimalWASMModule[4:8]) {
		return fmt.Errorf("wasm module version mismatch")
	}
	return nil
}

// RunWASMModule compiles, instantiates, and calls Peggy's deterministic
// Fibonacci WASM module with wazero.
// Intent: POC16 should now distinguish real runtime execution from the older
// header-only placeholder, while keeping arbitrary untrusted module loading out
// of scope for this POC. Source: DI-kimim; DI-sivis
func RunWASMModule(ctx context.Context, moduleBytes []byte, inputValue uint64) (result WASMRunResult, err error) {
	if err := ValidateWASMModule(moduleBytes); err != nil {
		return WASMRunResult{}, err
	}
	wasmRuntime := wazero.NewRuntime(ctx)
	defer func() {
		if closeErr := wasmRuntime.Close(ctx); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	compiledModule, err := wasmRuntime.CompileModule(ctx, moduleBytes)
	if err != nil {
		return WASMRunResult{}, fmt.Errorf("compile wasm module: %w", err)
	}
	module, err := wasmRuntime.InstantiateModule(ctx, compiledModule, wazero.NewModuleConfig().WithName("poc16-peggy"))
	if err != nil {
		return WASMRunResult{}, fmt.Errorf("instantiate wasm module: %w", err)
	}
	exportedFunction := module.ExportedFunction(WASMExportName)
	if exportedFunction == nil {
		return WASMRunResult{}, fmt.Errorf("wasm export %q not found", WASMExportName)
	}
	values, err := exportedFunction.Call(ctx, inputValue)
	if err != nil {
		return WASMRunResult{}, fmt.Errorf("call wasm export %q: %w", WASMExportName, err)
	}
	if len(values) != 1 {
		return WASMRunResult{}, fmt.Errorf("wasm export %q returned %d values, want 1", WASMExportName, len(values))
	}
	return WASMRunResult{ExportName: WASMExportName, InputValue: inputValue, ExportValue: values[0]}, nil
}

// PromiseFields returns the common relationship_v1 payload used by the
// heterogeneous runtime-adapter agents. The caller still signs and routes the
// envelope through the normal POC16 app/kernel path.
func PromiseFields(fromAgent, toAgent, promiseAbout, promiseText string) map[string]string {
	return map[string]string{
		"act":           "promise",
		"from":          fromAgent,
		"to":            toAgent,
		"turn":          "startup",
		"promise":       promiseText,
		"reason":        "heterogeneous runtime adapter events expressed as a local promise",
		"promise_about": promiseAbout,
	}
}

// StdioCBORRequest is Victor's adapter-to-worker request carried as a local
// length-prefixed CBOR frame. It is subprocess plumbing, not a PromiseGrid wire
// protocol and not an RPC command surface.
type StdioCBORRequest struct {
	Type string
	From string
	To   string
}

// StdioCBOREnvelope is the worker-to-adapter outbound promise envelope. The
// envelope bytes are carried as a CBOR byte string so stdout never hex-encodes
// the signed PromiseGrid message.
type StdioCBOREnvelope struct {
	Type          string
	From          string
	To            string
	Protocol      string
	EnvelopeBytes []byte
}

// StdioCBORAck is the adapter-to-worker acknowledgement envelope after the
// adapter receives peer events through the local kernel.
type StdioCBORAck struct {
	Type          string
	EnvelopeBytes []byte
}

// StdioCBOREvent is the worker's final stdout event after it parses the ACK
// envelope locally.
type StdioCBOREvent struct {
	Type     string
	Outcome  string
	ExactCID string
}

// WriteCBORFrame writes one bounded length-prefixed CBOR control frame.
// Intent: Victor's subprocess interface should carry binary CBOR records and
// exact envelope byte strings, replacing the previous JSON-plus-hex shim.
// Source: DI-kimim
func WriteCBORFrame(writer io.Writer, payloadBytes []byte) error {
	if len(payloadBytes) == 0 || len(payloadBytes) > maxStdioCBORFrameSize {
		return fmt.Errorf("invalid stdio cbor frame length: %d", len(payloadBytes))
	}
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(payloadBytes)))
	if err := writeAll(writer, header); err != nil {
		return err
	}
	return writeAll(writer, payloadBytes)
}

// ReadCBORFrame reads one bounded length-prefixed CBOR control frame.
func ReadCBORFrame(reader io.Reader) ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(header)
	if length == 0 || length > maxStdioCBORFrameSize {
		return nil, fmt.Errorf("invalid stdio cbor frame length: %d", length)
	}
	payloadBytes := make([]byte, int(length))
	if _, err := io.ReadFull(reader, payloadBytes); err != nil {
		return nil, err
	}
	return payloadBytes, nil
}

// MarshalStdioCBORRequest encodes a request as CBOR array
// [type, from, to].
func MarshalStdioCBORRequest(message StdioCBORRequest) ([]byte, error) {
	if message.Type == "" || message.From == "" || message.To == "" {
		return nil, fmt.Errorf("stdio cbor request is incomplete")
	}
	return marshalStdioCBORArray(
		stdioCBORString(message.Type),
		stdioCBORString(message.From),
		stdioCBORString(message.To),
	)
}

// ParseStdioCBORRequest decodes a request from CBOR array [type, from, to].
func ParseStdioCBORRequest(frameBytes []byte) (StdioCBORRequest, error) {
	reader, err := newStdioCBORArrayReader(frameBytes, 3)
	if err != nil {
		return StdioCBORRequest{}, err
	}
	message := StdioCBORRequest{}
	if message.Type, err = reader.readString(); err != nil {
		return StdioCBORRequest{}, err
	}
	if message.From, err = reader.readString(); err != nil {
		return StdioCBORRequest{}, err
	}
	if message.To, err = reader.readString(); err != nil {
		return StdioCBORRequest{}, err
	}
	if err := reader.finish(); err != nil {
		return StdioCBORRequest{}, err
	}
	if message.Type == "" || message.From == "" || message.To == "" {
		return StdioCBORRequest{}, fmt.Errorf("stdio cbor request is incomplete")
	}
	return message, nil
}

// MarshalStdioCBOREnvelope encodes an outbound envelope as CBOR array
// [type, from, to, protocol, envelope_bytes].
func MarshalStdioCBOREnvelope(message StdioCBOREnvelope) ([]byte, error) {
	if message.Type == "" || message.From == "" || message.To == "" || message.Protocol == "" || len(message.EnvelopeBytes) == 0 {
		return nil, fmt.Errorf("stdio cbor envelope is incomplete")
	}
	return marshalStdioCBORArray(
		stdioCBORString(message.Type),
		stdioCBORString(message.From),
		stdioCBORString(message.To),
		stdioCBORString(message.Protocol),
		stdioCBORBytes(message.EnvelopeBytes),
	)
}

// ParseStdioCBOREnvelope decodes an outbound envelope from CBOR array
// [type, from, to, protocol, envelope_bytes].
func ParseStdioCBOREnvelope(frameBytes []byte) (StdioCBOREnvelope, error) {
	reader, err := newStdioCBORArrayReader(frameBytes, 5)
	if err != nil {
		return StdioCBOREnvelope{}, err
	}
	message := StdioCBOREnvelope{}
	if message.Type, err = reader.readString(); err != nil {
		return StdioCBOREnvelope{}, err
	}
	if message.From, err = reader.readString(); err != nil {
		return StdioCBOREnvelope{}, err
	}
	if message.To, err = reader.readString(); err != nil {
		return StdioCBOREnvelope{}, err
	}
	if message.Protocol, err = reader.readString(); err != nil {
		return StdioCBOREnvelope{}, err
	}
	if message.EnvelopeBytes, err = reader.readBytes(); err != nil {
		return StdioCBOREnvelope{}, err
	}
	if err := reader.finish(); err != nil {
		return StdioCBOREnvelope{}, err
	}
	if message.Type == "" || message.From == "" || message.To == "" || message.Protocol == "" || len(message.EnvelopeBytes) == 0 {
		return StdioCBOREnvelope{}, fmt.Errorf("stdio cbor envelope is incomplete")
	}
	return message, nil
}

// MarshalStdioCBORAck encodes an ACK as CBOR array [type, envelope_bytes].
func MarshalStdioCBORAck(message StdioCBORAck) ([]byte, error) {
	if message.Type == "" || len(message.EnvelopeBytes) == 0 {
		return nil, fmt.Errorf("stdio cbor ack is incomplete")
	}
	return marshalStdioCBORArray(
		stdioCBORString(message.Type),
		stdioCBORBytes(message.EnvelopeBytes),
	)
}

// ParseStdioCBORAck decodes an ACK from CBOR array [type, envelope_bytes].
func ParseStdioCBORAck(frameBytes []byte) (StdioCBORAck, error) {
	reader, err := newStdioCBORArrayReader(frameBytes, 2)
	if err != nil {
		return StdioCBORAck{}, err
	}
	message := StdioCBORAck{}
	if message.Type, err = reader.readString(); err != nil {
		return StdioCBORAck{}, err
	}
	if message.EnvelopeBytes, err = reader.readBytes(); err != nil {
		return StdioCBORAck{}, err
	}
	if err := reader.finish(); err != nil {
		return StdioCBORAck{}, err
	}
	if message.Type == "" || len(message.EnvelopeBytes) == 0 {
		return StdioCBORAck{}, fmt.Errorf("stdio cbor ack is incomplete")
	}
	return message, nil
}

// MarshalStdioCBOREvent encodes a worker event as CBOR array
// [type, outcome, exact_cid].
func MarshalStdioCBOREvent(message StdioCBOREvent) ([]byte, error) {
	if message.Type == "" || message.Outcome == "" || message.ExactCID == "" {
		return nil, fmt.Errorf("stdio cbor event is incomplete")
	}
	return marshalStdioCBORArray(
		stdioCBORString(message.Type),
		stdioCBORString(message.Outcome),
		stdioCBORString(message.ExactCID),
	)
}

// ParseStdioCBOREvent decodes a worker event from CBOR array
// [type, outcome, exact_cid].
func ParseStdioCBOREvent(frameBytes []byte) (StdioCBOREvent, error) {
	reader, err := newStdioCBORArrayReader(frameBytes, 3)
	if err != nil {
		return StdioCBOREvent{}, err
	}
	message := StdioCBOREvent{}
	if message.Type, err = reader.readString(); err != nil {
		return StdioCBOREvent{}, err
	}
	if message.Outcome, err = reader.readString(); err != nil {
		return StdioCBOREvent{}, err
	}
	if message.ExactCID, err = reader.readString(); err != nil {
		return StdioCBOREvent{}, err
	}
	if err := reader.finish(); err != nil {
		return StdioCBOREvent{}, err
	}
	if message.Type == "" || message.Outcome == "" || message.ExactCID == "" {
		return StdioCBOREvent{}, fmt.Errorf("stdio cbor event is incomplete")
	}
	return message, nil
}

type stdioCBORItem struct {
	text    string
	bytes   []byte
	isBytes bool
}

func stdioCBORString(value string) stdioCBORItem {
	return stdioCBORItem{text: value}
}

func stdioCBORBytes(value []byte) stdioCBORItem {
	copied := append([]byte(nil), value...)
	return stdioCBORItem{bytes: copied, isBytes: true}
}

type stdioCBORWriter struct {
	buffer bytes.Buffer
}

// marshalStdioCBORArray writes only the small deterministic CBOR subset Victor
// needs: definite-length arrays containing text strings and byte strings.
// Intent: Keep stdio framing binary and exact-byte preserving without promoting
// this local adapter contract into a general-purpose PromiseGrid wire schema.
// Source: DI-kimim
func marshalStdioCBORArray(items ...stdioCBORItem) ([]byte, error) {
	writer := &stdioCBORWriter{}
	if err := writer.writeArrayHeader(len(items)); err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.isBytes {
			if err := writer.writeBytes(item.bytes); err != nil {
				return nil, err
			}
			continue
		}
		if err := writer.writeString(item.text); err != nil {
			return nil, err
		}
	}
	return writer.buffer.Bytes(), nil
}

func (writer *stdioCBORWriter) writeTypeAndLength(major byte, length uint64) error {
	prefix := major << 5
	switch {
	case length < 24:
		return writer.buffer.WriteByte(prefix | byte(length))
	case length <= 0xff:
		if err := writer.buffer.WriteByte(prefix | 24); err != nil {
			return err
		}
		return writer.buffer.WriteByte(byte(length))
	case length <= 0xffff:
		if err := writer.buffer.WriteByte(prefix | 25); err != nil {
			return err
		}
		return binary.Write(&writer.buffer, binary.BigEndian, uint16(length))
	case length <= 0xffffffff:
		if err := writer.buffer.WriteByte(prefix | 26); err != nil {
			return err
		}
		return binary.Write(&writer.buffer, binary.BigEndian, uint32(length))
	default:
		if err := writer.buffer.WriteByte(prefix | 27); err != nil {
			return err
		}
		return binary.Write(&writer.buffer, binary.BigEndian, length)
	}
}

func (writer *stdioCBORWriter) writeArrayHeader(length int) error {
	return writer.writeTypeAndLength(4, uint64(length))
}

func (writer *stdioCBORWriter) writeBytes(value []byte) error {
	if err := writer.writeTypeAndLength(2, uint64(len(value))); err != nil {
		return err
	}
	_, err := writer.buffer.Write(value)
	return err
}

func (writer *stdioCBORWriter) writeString(value string) error {
	if !utf8.ValidString(value) {
		return fmt.Errorf("invalid stdio cbor utf8 string")
	}
	if err := writer.writeTypeAndLength(3, uint64(len(value))); err != nil {
		return err
	}
	_, err := writer.buffer.WriteString(value)
	return err
}

type stdioCBORReader struct {
	data   []byte
	offset int
}

// newStdioCBORArrayReader rejects unexpected frame shapes early so the worker
// and adapter fail closed on malformed local subprocess frames.
func newStdioCBORArrayReader(frameBytes []byte, expectedLength uint64) (*stdioCBORReader, error) {
	reader := &stdioCBORReader{data: frameBytes}
	length, err := reader.readTypeAndLength(4)
	if err != nil {
		return nil, err
	}
	if length != expectedLength {
		return nil, fmt.Errorf("stdio cbor array length = %d, want %d", length, expectedLength)
	}
	return reader, nil
}

func (reader *stdioCBORReader) readByte() (byte, error) {
	if reader.offset >= len(reader.data) {
		return 0, fmt.Errorf("unexpected end of stdio cbor data")
	}
	value := reader.data[reader.offset]
	reader.offset++
	return value, nil
}

func (reader *stdioCBORReader) readTypeAndLength(expectedMajor byte) (uint64, error) {
	initial, err := reader.readByte()
	if err != nil {
		return 0, err
	}
	major := initial >> 5
	if major != expectedMajor {
		return 0, fmt.Errorf("expected stdio cbor major %d, got %d", expectedMajor, major)
	}
	additional := initial & 0x1f
	switch additional {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23:
		return uint64(additional), nil
	case 24:
		value, readErr := reader.readByte()
		return uint64(value), readErr
	case 25:
		if reader.offset+2 > len(reader.data) {
			return 0, fmt.Errorf("truncated stdio cbor uint16")
		}
		value := binary.BigEndian.Uint16(reader.data[reader.offset : reader.offset+2])
		reader.offset += 2
		return uint64(value), nil
	case 26:
		if reader.offset+4 > len(reader.data) {
			return 0, fmt.Errorf("truncated stdio cbor uint32")
		}
		value := binary.BigEndian.Uint32(reader.data[reader.offset : reader.offset+4])
		reader.offset += 4
		return uint64(value), nil
	case 27:
		if reader.offset+8 > len(reader.data) {
			return 0, fmt.Errorf("truncated stdio cbor uint64")
		}
		value := binary.BigEndian.Uint64(reader.data[reader.offset : reader.offset+8])
		reader.offset += 8
		return value, nil
	default:
		return 0, fmt.Errorf("unsupported stdio cbor additional information %d", additional)
	}
}

func (reader *stdioCBORReader) readBytes() ([]byte, error) {
	length, err := reader.readTypeAndLength(2)
	if err != nil {
		return nil, err
	}
	if reader.offset+int(length) > len(reader.data) {
		return nil, fmt.Errorf("truncated stdio cbor byte string")
	}
	value := make([]byte, int(length))
	copy(value, reader.data[reader.offset:reader.offset+int(length)])
	reader.offset += int(length)
	return value, nil
}

func (reader *stdioCBORReader) readString() (string, error) {
	length, err := reader.readTypeAndLength(3)
	if err != nil {
		return "", err
	}
	if reader.offset+int(length) > len(reader.data) {
		return "", fmt.Errorf("truncated stdio cbor string")
	}
	value := string(reader.data[reader.offset : reader.offset+int(length)])
	reader.offset += int(length)
	if !utf8.ValidString(value) {
		return "", fmt.Errorf("invalid stdio cbor utf8 string")
	}
	return value, nil
}

func (reader *stdioCBORReader) finish() error {
	if reader.offset != len(reader.data) {
		return fmt.Errorf("trailing stdio cbor bytes: %d", len(reader.data)-reader.offset)
	}
	return nil
}

// writeAll handles short writes on subprocess pipes so a CBOR frame is either
// completely written or reported as a local adapter failure.
func writeAll(writer io.Writer, payloadBytes []byte) error {
	for len(payloadBytes) > 0 {
		written, writeErr := writer.Write(payloadBytes)
		if written > 0 {
			payloadBytes = payloadBytes[written:]
		}
		if writeErr != nil {
			return writeErr
		}
		if written == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}
