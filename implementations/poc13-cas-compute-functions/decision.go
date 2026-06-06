package poc13

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Decider supplies local autonomy text for one agent turn.
type Decider interface {
	Decide(ctx context.Context, cfg Config, agent AgentConfig, prompt string) (DecisionResult, error)
}

// DecisionResult records whether autonomy came from a live provider or a local
// fallback.
type DecisionResult struct {
	Mode string
	Text string
}

// LiveOrScriptedDecider uses a provider only when configured and credentialed;
// otherwise it returns an explicit local fallback.
// Intent: POC13 should be live-LLM-capable without making tests, Docker builds,
// or config files depend on storing provider credentials. Source: DI-notig
type LiveOrScriptedDecider struct{}

// Decide returns provider text when possible and clear fallback text otherwise.
func (decider LiveOrScriptedDecider) Decide(ctx context.Context, cfg Config, agent AgentConfig, prompt string) (DecisionResult, error) {
	apiKey := os.Getenv(cfg.APIKeyEnv)
	if !cfg.LiveDecisions || apiKey == "" {
		return DecisionResult{Mode: "scripted", Text: "local fallback: promise only what I can verify and keep"}, nil
	}
	requestBody := map[string]any{
		"model": cfg.AgentModel,
		"input": []map[string]string{{
			"role":    "system",
			"content": "You are a PromiseGrid POC13 agent. Respond with one concise local promise judgment, no commands, no permissions, no authority claims.",
		}, {
			"role":    "user",
			"content": agent.Persona + "\n" + prompt,
		}},
		"reasoning": map[string]string{"effort": cfg.ReasoningEffort},
	}
	if cfg.ServiceTier != "" {
		requestBody["service_tier"] = cfg.ServiceTier
	}
	bodyBytes, marshalErr := json.Marshal(requestBody)
	if marshalErr != nil {
		return DecisionResult{}, marshalErr
	}
	request, requestErr := http.NewRequestWithContext(ctx, http.MethodPost, cfg.ProviderBaseURL, bytes.NewReader(bodyBytes))
	if requestErr != nil {
		return DecisionResult{}, requestErr
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: cfg.Timeout()}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		return DecisionResult{}, responseErr
	}
	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc13: close provider response: %v\n", closeErr)
		}
	}()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return DecisionResult{}, fmt.Errorf("provider status %d", response.StatusCode)
	}
	var responseBody struct {
		OutputText string `json:"output_text"`
	}
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return DecisionResult{}, err
	}
	if responseBody.OutputText == "" {
		responseBody.OutputText = "provider returned no output_text"
	}
	return DecisionResult{Mode: "live", Text: responseBody.OutputText}, nil
}
