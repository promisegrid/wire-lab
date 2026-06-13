package decision

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// LiveClient is a small POC-local OpenAI-compatible Responses API client.
// Intent: POC14 must not import or shell out to ga-runner; this client is
// deliberately narrow and exists only for the autonomous-agent POC. Source: DI-timah
type LiveClient struct {
	BaseURL         string
	APIKeyEnv       string
	AgentModel      string
	MonitorModel    string
	ReasoningEffort string
	ServiceTier     string
	HTTPClient      *http.Client
}

type liveDecisionResponse struct {
	Act     string              `json:"act"`
	Target  string              `json:"target"`
	Promise string              `json:"promise"`
	Reason  string              `json:"reason"`
	Fields  []liveDecisionField `json:"fields"`
}

type liveDecisionField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// promiseDecision converts the strict provider shape back into the runtime's
// map form. The wire payload still uses pCID-defined string fields; the array
// exists only because strict Structured Outputs require closed object schemas.
// Intent: Use provider schema enforcement without changing the POC14 protocol
// payload vocabulary or allowing arbitrary provider-shaped objects.
// Source: DI-timah
func (response liveDecisionResponse) promiseDecision() PromiseDecision {
	fields := make(map[string]any, len(response.Fields))
	for _, field := range response.Fields {
		key := strings.TrimSpace(field.Key)
		if key != "" {
			fields[key] = field.Value
		}
	}
	return PromiseDecision{
		Act:     response.Act,
		Target:  response.Target,
		Promise: response.Promise,
		Reason:  response.Reason,
		Fields:  fields,
	}
}

// NewLiveClient constructs a live decider with an explicit timeout.
func NewLiveClient(baseURL, apiKeyEnv, agentModel, monitorModel, reasoningEffort, serviceTier string, timeout time.Duration) LiveClient {
	return LiveClient{
		BaseURL:         baseURL,
		APIKeyEnv:       apiKeyEnv,
		AgentModel:      agentModel,
		MonitorModel:    monitorModel,
		ReasoningEffort: reasoningEffort,
		ServiceTier:     serviceTier,
		HTTPClient:      &http.Client{Timeout: timeout},
	}
}

// Decide asks the live provider for one local agent promise decision.
func (client LiveClient) Decide(ctx context.Context, observation Observation) (PromiseDecision, error) {
	prompt, promptErr := Prompt(observation)
	if promptErr != nil {
		return PromiseDecision{}, promptErr
	}
	text, callErr := client.callResponses(ctx, client.AgentModel, agentSystemPrompt(), prompt, promiseDecisionTextFormat(observation))
	if callErr != nil {
		return PromiseDecision{}, callErr
	}
	var liveDecision liveDecisionResponse
	if err := json.Unmarshal([]byte(extractJSONObject(text)), &liveDecision); err != nil {
		return PromiseDecision{}, fmt.Errorf("decode live decision JSON: %w: %s", err, text)
	}
	return liveDecision.promiseDecision(), nil
}

// Evaluate asks the live provider for an observer-only monitor report.
func (client LiveClient) Evaluate(ctx context.Context, events []Event) (MonitorReport, error) {
	encodedEvents, marshalErr := json.MarshalIndent(events, "", "  ")
	if marshalErr != nil {
		return MonitorReport{}, marshalErr
	}
	userPrompt := "Evaluate this POC14 run as an observer only. Return exactly one JSON object matching " +
		`{"promise_theory_fit":0,"autonomy":0,"protocol_validity":0,"local_trust_correctness":0,"imposition_avoidance":0,"summary":"","concerns":[]}` +
		". Score each integer from 0 to 5. Do not issue commands to agents.\n\nEvents:\n" + string(encodedEvents)
	text, callErr := client.callResponses(ctx, client.MonitorModel, monitorSystemPrompt(), userPrompt, monitorReportTextFormat())
	if callErr != nil {
		return MonitorReport{}, callErr
	}
	var report MonitorReport
	if err := json.Unmarshal([]byte(extractJSONObject(text)), &report); err != nil {
		return MonitorReport{}, fmt.Errorf("decode live monitor JSON: %w: %s", err, text)
	}
	return report, nil
}

// callResponses sends the compact POC-local Responses request. The optional
// text format asks the provider to enforce the same JSON object shape that the
// Go validator checks again after receipt.
// Intent: Reduce malformed live decisions without letting provider schema
// enforcement become the source of PromiseGrid semantics. Source: DI-timah
func (client LiveClient) callResponses(ctx context.Context, model, systemPrompt, userPrompt string, textFormat map[string]any) (string, error) {
	apiKey, apiKeyErr := client.apiKey()
	if apiKeyErr != nil {
		return "", apiKeyErr
	}
	requestBody := map[string]any{
		"model":        model,
		"service_tier": client.ServiceTier,
		"input": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
	}
	if client.ReasoningEffort != "" {
		requestBody["reasoning"] = map[string]string{"effort": client.ReasoningEffort}
	}
	if textFormat != nil {
		requestBody["text"] = map[string]any{"format": textFormat}
	}
	encodedBody, marshalErr := json.Marshal(requestBody)
	if marshalErr != nil {
		return "", marshalErr
	}
	request, requestErr := http.NewRequestWithContext(ctx, http.MethodPost, client.BaseURL, bytes.NewReader(encodedBody))
	if requestErr != nil {
		return "", requestErr
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")
	response, responseErr := client.HTTPClient.Do(request)
	if responseErr != nil {
		return "", responseErr
	}
	defer func() {
		_, copyErr := io.Copy(io.Discard, response.Body)
		if copyErr != nil {
			return
		}
		closeErr := response.Body.Close()
		if closeErr != nil {
			return
		}
	}()
	responseBytes, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return "", readErr
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("provider status %d: %s", response.StatusCode, string(responseBytes))
	}
	outputText, parseErr := parseResponseText(responseBytes)
	if parseErr != nil {
		return "", parseErr
	}
	return outputText, nil
}

// promiseDecisionTextFormat is the provider-side schema for one live agent
// decision; it intentionally permits only the single top-level promise action
// and constrains targets to names visible in this turn's local observation.
// Intent: Keep repairs, refusals, observations, economics, and candidate-link
// discovery inside promise content while preventing arbitrary target strings.
// Source: DI-timah; DI-bikit
func promiseDecisionTextFormat(observation Observation) map[string]any {
	targetSchema := map[string]any{"type": "string"}
	if targets := structuredTargetNames(observation); len(targets) > 0 {
		targetSchema["enum"] = targets
	}
	return map[string]any{
		"type":   "json_schema",
		"name":   "poc14_promise_decision",
		"strict": true,
		"schema": map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"required":             []string{"act", "target", "promise", "reason", "fields"},
			"properties": map[string]any{
				"act":     map[string]any{"type": "string", "enum": []string{ActPromise}},
				"target":  targetSchema,
				"promise": map[string]any{"type": "string"},
				"reason":  map[string]any{"type": "string"},
				"fields": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"required":             []string{"key", "value"},
						"properties": map[string]any{
							"key":   map[string]any{"type": "string"},
							"value": map[string]any{"type": "string"},
						},
					},
				},
			},
		},
	}
}

// structuredTargetNames returns the provider-visible target enum for one turn.
// Intent: The live model can choose one visible direct peer, or one candidate
// peer only when the later Go validator sees a link-discovery promise payload.
// Source: DI-timah
func structuredTargetNames(observation Observation) []string {
	seenTargets := make(map[string]bool)
	var targets []string
	for _, target := range append(append([]string{}, observation.DirectPeers...), observation.CandidatePeers...) {
		target = strings.TrimSpace(target)
		if target == "" || target == observation.AgentName || seenTargets[target] {
			continue
		}
		seenTargets[target] = true
		targets = append(targets, target)
	}
	return targets
}

// monitorReportTextFormat is the provider-side schema for the observer-only
// monitor report.
// Intent: Keep monitor output structured and comparable while preserving the
// monitor as local events, not an authority over agents. Source: DI-timah
func monitorReportTextFormat() map[string]any {
	return map[string]any{
		"type":   "json_schema",
		"name":   "poc14_monitor_report",
		"strict": true,
		"schema": map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"required":             []string{"promise_theory_fit", "autonomy", "protocol_validity", "local_trust_correctness", "imposition_avoidance", "summary", "concerns"},
			"properties": map[string]any{
				"promise_theory_fit":      map[string]any{"type": "integer", "minimum": 0, "maximum": 5},
				"autonomy":                map[string]any{"type": "integer", "minimum": 0, "maximum": 5},
				"protocol_validity":       map[string]any{"type": "integer", "minimum": 0, "maximum": 5},
				"local_trust_correctness": map[string]any{"type": "integer", "minimum": 0, "maximum": 5},
				"imposition_avoidance":    map[string]any{"type": "integer", "minimum": 0, "maximum": 5},
				"summary":                 map[string]any{"type": "string"},
				"concerns":                map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			},
		},
	}
}

func (client LiveClient) apiKey() (string, error) {
	apiKey := strings.TrimSpace(os.Getenv(client.APIKeyEnv))
	if apiKey != "" {
		return apiKey, nil
	}
	fileEnv := client.APIKeyEnv + "_FILE"
	filePath := strings.TrimSpace(os.Getenv(fileEnv))
	if filePath == "" {
		return "", fmt.Errorf("environment variable %s or %s is required", client.APIKeyEnv, fileEnv)
	}
	// Intent: Compose secret files avoid rendering provider key values in
	// `docker compose config`, while preserving the same config-level API-key
	// environment name. Source: DI-timah
	fileBytes, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return "", readErr
	}
	apiKey = strings.TrimSpace(string(fileBytes))
	if apiKey == "" {
		return "", fmt.Errorf("API key file %s is empty", filePath)
	}
	return apiKey, nil
}

func parseResponseText(responseBytes []byte) (string, error) {
	var response struct {
		OutputText string `json:"output_text"`
		Output     []struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return "", err
	}
	if strings.TrimSpace(response.OutputText) != "" {
		return response.OutputText, nil
	}
	var builder strings.Builder
	for _, output := range response.Output {
		for _, content := range output.Content {
			builder.WriteString(content.Text)
		}
	}
	text := strings.TrimSpace(builder.String())
	if text == "" {
		return "", fmt.Errorf("provider response did not include output text")
	}
	return text, nil
}

func extractJSONObject(text string) string {
	trimmed := strings.TrimSpace(text)
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end >= start {
		return trimmed[start : end+1]
	}
	return trimmed
}

func agentSystemPrompt() string {
	return "You are one autonomous PromiseGrid agent. The only valid top-level act is promise. You may promise only your own behavior. Refusals, repairs, observations, economics, and link preferences are promise content, not action kinds. All trust is local. Do not use permission, authorization, conformance, enforcement, or contract language as external authority. Return only JSON."
}

func monitorSystemPrompt() string {
	return "You are an observer-only PromiseGrid monitor. Evaluate logs for Promise Theory fit, autonomy, protocol validity, local trust correctness, and imposition avoidance. You are not an authority and must not control agents. Return only JSON."
}
