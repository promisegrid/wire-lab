package decision

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAPIKeyReadsDirectEnvironment(t *testing.T) {
	t.Setenv("POC12_TEST_API_KEY", "direct-test-key")
	client := NewLiveClient("https://example.invalid/v1/responses", "POC12_TEST_API_KEY", "agent", "monitor", "medium", "flex", time.Second)
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
	t.Setenv("POC12_TEST_API_KEY_FILE", secretPath)
	client := NewLiveClient("https://example.invalid/v1/responses", "POC12_TEST_API_KEY", "agent", "monitor", "medium", "flex", time.Second)
	apiKey, err := client.apiKey()
	if err != nil {
		t.Fatalf("api key: %v", err)
	}
	if apiKey != "file-test-key" {
		t.Fatalf("api key = %q, want file-test-key", apiKey)
	}
}

func TestAPIKeyRejectsMissingSources(t *testing.T) {
	client := NewLiveClient("https://example.invalid/v1/responses", "POC12_TEST_MISSING_API_KEY", "agent", "monitor", "medium", "flex", time.Second)
	if _, err := client.apiKey(); err == nil {
		t.Fatalf("missing API key sources should fail")
	}
}

func TestLiveClientRequestsStructuredDecisionOutput(t *testing.T) {
	t.Setenv("POC12_TEST_API_KEY", "test-key")
	var requestBody map[string]any
	client := NewLiveClient("https://example.invalid/v1/responses", "POC12_TEST_API_KEY", "agent", "monitor", "medium", "flex", time.Second)
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		bodyBytes, readErr := io.ReadAll(request.Body)
		if readErr != nil {
			return nil, readErr
		}
		if closeErr := request.Body.Close(); closeErr != nil {
			return nil, closeErr
		}
		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			return nil, err
		}
		responseText := `{"output_text":"{\"act\":\"promise\",\"target\":\"bob\",\"promise\":\"Alice promises one local exchange.\",\"reason\":\"test\",\"fields\":[{\"key\":\"promise_about\",\"value\":\"local_observation\"}]}"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(responseText)),
			Header:     make(http.Header),
			Request:    request,
		}, nil
	})}
	promiseDecision, err := client.Decide(context.Background(), Observation{AgentName: "alice", DirectPeers: []string{"bob"}, CandidatePeers: []string{"carol"}})
	if err != nil {
		t.Fatalf("decide: %v", err)
	}
	if promiseDecision.Act != ActPromise {
		t.Fatalf("decision act = %q, want promise", promiseDecision.Act)
	}
	if promiseDecision.Fields["promise_about"] != "local_observation" {
		t.Fatalf("decision fields not converted from strict live shape: %#v", promiseDecision.Fields)
	}
	textObject, ok := requestBody["text"].(map[string]any)
	if !ok {
		t.Fatalf("request missing text format: %#v", requestBody)
	}
	formatObject, ok := textObject["format"].(map[string]any)
	if !ok {
		t.Fatalf("request missing format object: %#v", textObject)
	}
	if formatObject["type"] != "json_schema" || formatObject["name"] != "poc12_promise_decision" {
		t.Fatalf("unexpected structured output format: %#v", formatObject)
	}
	schemaObject := formatObject["schema"].(map[string]any)
	propertiesObject := schemaObject["properties"].(map[string]any)
	fieldsObject := propertiesObject["fields"].(map[string]any)
	if fieldsObject["type"] != "array" {
		t.Fatalf("live fields should use a strict key/value array: %#v", fieldsObject)
	}
	targetObject := propertiesObject["target"].(map[string]any)
	targetEnum := targetObject["enum"].([]any)
	if len(targetEnum) != 2 || targetEnum[0] != "bob" || targetEnum[1] != "carol" {
		t.Fatalf("target enum = %#v, want bob and carol", targetEnum)
	}
}

// roundTripFunc lets live-client tests inspect provider requests without opening
// a test TCP listener. Intent: Keep tests deterministic and sandbox-compatible
// while still validating the provider request shape. Source: DI-timah
type roundTripFunc func(*http.Request) (*http.Response, error)

func (roundTripFuncValue roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTripFuncValue(request)
}
