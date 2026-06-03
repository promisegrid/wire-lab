package decision

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAPIKeyReadsDirectEnvironment(t *testing.T) {
	t.Setenv("POC11_TEST_API_KEY", "direct-test-key")
	client := NewLiveClient("https://example.invalid/v1/responses", "POC11_TEST_API_KEY", "agent", "monitor", "medium", "flex", time.Second)
	apiKey, err := client.apiKey()
	if err != nil {
		t.Fatalf("api key: %v", err)
	}
	if apiKey != "direct-test-key" {
		t.Fatalf("api key = %q, want direct-test-key", apiKey)
	}
}

func TestAPIKeyReadsFileEnvironment(t *testing.T) {
	secretPath := filepath.Join(t.TempDir(), "openai_api_key.txt")
	if err := os.WriteFile(secretPath, []byte("file-test-key\n"), 0o600); err != nil {
		t.Fatalf("write secret file: %v", err)
	}
	t.Setenv("POC11_TEST_API_KEY_FILE", secretPath)
	client := NewLiveClient("https://example.invalid/v1/responses", "POC11_TEST_API_KEY", "agent", "monitor", "medium", "flex", time.Second)
	apiKey, err := client.apiKey()
	if err != nil {
		t.Fatalf("api key: %v", err)
	}
	if apiKey != "file-test-key" {
		t.Fatalf("api key = %q, want file-test-key", apiKey)
	}
}

func TestAPIKeyRejectsMissingSources(t *testing.T) {
	client := NewLiveClient("https://example.invalid/v1/responses", "POC11_TEST_MISSING_API_KEY", "agent", "monitor", "medium", "flex", time.Second)
	if _, err := client.apiKey(); err == nil {
		t.Fatalf("missing API key sources should fail")
	}
}
