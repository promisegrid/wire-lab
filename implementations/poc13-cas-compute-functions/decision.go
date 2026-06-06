package poc13

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

// Promises reports whether the provider/local judgment contains a voluntary
// promise rather than a non-commitment.
// Intent: POC13 protocol behavior should be gated by the agent's local promise
// judgment, so LLM text can stop action instead of being passive decoration.
// Source: DI-fumol
func (decision DecisionResult) Promises() bool {
	lowerText := strings.ToLower(decision.Text)
	if !strings.Contains(lowerText, "promise") {
		return false
	}
	for _, negativePhrase := range []string{"do not promise", "cannot promise", "can't promise", "no promise", "non-commitment"} {
		if strings.Contains(lowerText, negativePhrase) {
			return false
		}
	}
	return true
}

// ProviderResponse models only the stable text-bearing parts of the Responses
// API shape POC13 needs for local promise evidence.
// Intent: Live runs should record the provider's actual local promise judgment
// instead of a placeholder when nested output text is present. Source:
// DI-lasuh
type ProviderResponse struct {
	OutputText string           `json:"output_text"`
	Output     []ProviderOutput `json:"output"`
}

// ProviderOutput is one Responses API output item.
type ProviderOutput struct {
	Content []ProviderContent `json:"content"`
}

// ProviderContent is one text-bearing content item inside a provider output.
type ProviderContent struct {
	Text string `json:"text"`
}

// ResponseText returns the first non-empty text field found in the provider
// response, checking both the top-level convenience field and nested content.
func (response ProviderResponse) ResponseText() string {
	if text := strings.TrimSpace(response.OutputText); text != "" {
		return text
	}
	for _, output := range response.Output {
		for _, content := range output.Content {
			if text := strings.TrimSpace(content.Text); text != "" {
				return text
			}
		}
	}
	return ""
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
	var responseBody ProviderResponse
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return DecisionResult{}, err
	}
	responseText := responseBody.ResponseText()
	if responseText == "" {
		responseText = "provider returned no output_text"
	}
	return DecisionResult{Mode: "live", Text: responseText}, nil
}
