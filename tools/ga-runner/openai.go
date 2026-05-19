package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAIProvider struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// Generate calls the OpenAI-compatible Responses API and normalizes the result
// into the provider-neutral shape used by GA state and result writers.
//
// Intent: Keep the first provider-backed GA implementation concrete and
// checkpointable without leaking OpenAI response details into score/generate
// orchestration. Source: DI-gijom
func (provider OpenAIProvider) Generate(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
	if provider.APIKey == "" {
		return ProviderResponse{}, fmt.Errorf("missing OpenAI API key")
	}
	baseURL := provider.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/responses"
	}
	client := provider.Client
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Minute}
	}
	body := openAIRequest{
		Model:           request.APIModel,
		Input:           request.Prompt,
		Instructions:    request.Instructions,
		MaxOutputTokens: request.MaxOutputTokens,
	}
	if request.ReasoningEffort != "" {
		body.Reasoning = &openAIReasoning{Effort: request.ReasoningEffort}
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return ProviderResponse{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(payload))
	if err != nil {
		return ProviderResponse{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+provider.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return ProviderResponse{}, err
	}
	responseBytes, err := io.ReadAll(httpResp.Body)
	closeErr := httpResp.Body.Close()
	if err != nil {
		return ProviderResponse{}, err
	}
	if closeErr != nil {
		return ProviderResponse{}, closeErr
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return ProviderResponse{}, fmt.Errorf("openai status %d: %s", httpResp.StatusCode, strings.TrimSpace(string(responseBytes)))
	}
	var parsed openAIResponse
	if err := json.Unmarshal(responseBytes, &parsed); err != nil {
		return ProviderResponse{}, err
	}
	text := parsed.OutputText
	if text == "" {
		text = parsed.JoinOutputText()
	}
	if strings.TrimSpace(text) == "" {
		return ProviderResponse{}, fmt.Errorf("openai response contained no output text")
	}
	usageJSON := ""
	if len(parsed.Usage) > 0 {
		usageBytes, err := json.Marshal(parsed.Usage)
		if err != nil {
			return ProviderResponse{}, err
		}
		usageJSON = string(usageBytes)
	}
	if parsed.Status != "" && parsed.Status != "completed" {
		return ProviderResponse{}, fmt.Errorf("openai response status %q with output text length %d", parsed.Status, len(text))
	}
	return ProviderResponse{
		Text:       strings.TrimSpace(text) + "\n",
		RequestID:  httpResp.Header.Get("x-request-id"),
		ResponseID: parsed.ID,
		UsageJSON:  usageJSON,
	}, nil
}

type openAIRequest struct {
	Model           string           `json:"model"`
	Input           string           `json:"input"`
	Instructions    string           `json:"instructions,omitempty"`
	MaxOutputTokens int              `json:"max_output_tokens,omitempty"`
	Reasoning       *openAIReasoning `json:"reasoning,omitempty"`
}

type openAIReasoning struct {
	Effort string `json:"effort"`
}

type openAIResponse struct {
	ID         string                 `json:"id"`
	Status     string                 `json:"status"`
	OutputText string                 `json:"output_text"`
	Output     []openAIOutputItem     `json:"output"`
	Usage      map[string]interface{} `json:"usage"`
}

type openAIOutputItem struct {
	Type    string              `json:"type"`
	Content []openAIContentPart `json:"content"`
}

type openAIContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (response openAIResponse) JoinOutputText() string {
	var parts []string
	for _, item := range response.Output {
		if item.Type != "message" {
			continue
		}
		for _, content := range item.Content {
			if content.Type == "output_text" && content.Text != "" {
				parts = append(parts, content.Text)
			}
		}
	}
	return strings.Join(parts, "\n")
}
