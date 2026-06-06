package poc13

import (
	"bytes"
	"testing"
)

func TestFrameRoundTrip(t *testing.T) {
	var buffer bytes.Buffer
	want := []byte{0xd8, 0x1a, 0x67, 0x72, 0x69, 0x64}
	if err := (FrameWriter{writer: &buffer}).WriteFrame(want); err != nil {
		t.Fatalf("write frame: %v", err)
	}
	got, err := (FrameReader{reader: &buffer}).ReadFrame()
	if err != nil {
		t.Fatalf("read frame: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("frame = %x, want %x", got, want)
	}
}

func TestExecuteFunctionUsesPayloadInput(t *testing.T) {
	resultBytes, err := ExecuteFunction([]byte("poc13 function: fibonacci(n) v1"), []byte("n=10"), sampleContextBytes())
	if err != nil {
		t.Fatalf("execute function: %v", err)
	}
	want := "fibonacci(10)=55;context_cid=" + ContentCID(sampleContextBytes())
	if string(resultBytes) != want {
		t.Fatalf("result = %q, want %q", string(resultBytes), want)
	}
}

func TestExecuteFunctionRejectsUnsupportedSource(t *testing.T) {
	if _, err := ExecuteFunction([]byte("delete everything"), []byte("n=10"), sampleContextBytes()); err == nil {
		t.Fatalf("unsupported function source should fail")
	}
}
