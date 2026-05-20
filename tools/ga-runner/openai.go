package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type OpenAIProvider struct {
	APIKey      string
	BaseURL     string
	Client      *http.Client
	RetryPolicy ProviderRetryPolicy
	// DebugWriter receives raw request and response diagnostics for provider
	// calls. Authorization headers are intentionally excluded from these logs.
	DebugWriter io.Writer
}

// ProviderRetryPolicy bounds provider retries for long-running unattended GA
// runs. The zero value selects conservative sync canary defaults.
//
// Intent: Keep retry timing explicit, finite, and testable so Flex scarcity and
// slow provider calls do not require interactive babysitting. Source: DI-mopob;
// DI-juzus
type ProviderRetryPolicy struct {
	MaxAttempts    int
	MaxElapsed     time.Duration
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// Generate calls the OpenAI-compatible Responses API and normalizes the result
// into the provider-neutral shape used by GA state and result writers.
//
// Intent: Keep the first provider-backed GA implementation concrete and
// checkpointable without leaking OpenAI response details into score/generate
// orchestration. The wrapper now applies bounded Flex retries before returning a
// checkpointable provider error. Source: DI-gijom; DI-mopob
func (provider OpenAIProvider) Generate(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
	policy := provider.RetryPolicy.withDefaults()
	var lastErr error
	startedAt := time.Now()
	attempts := 0
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		attempts = attempt
		response, err := provider.generateOnce(ctx, request, attempt)
		if err == nil {
			return response, nil
		}
		lastErr = err
		provider.debugf("attempt=%d event=error elapsed=%s error=%q", attempt, time.Since(startedAt).Round(time.Millisecond), err.Error())
		if !shouldRetryOpenAIError(ctx, err) || attempt == policy.MaxAttempts {
			break
		}
		backoff := policy.backoff(attempt)
		if policy.MaxElapsed > 0 && time.Since(startedAt)+backoff > policy.MaxElapsed {
			break
		}
		provider.debugf("attempt=%d event=retry_sleep duration=%s", attempt, backoff)
		if err := sleepWithContext(ctx, backoff); err != nil {
			return ProviderResponse{}, err
		}
	}
	if attempts == 1 {
		return ProviderResponse{}, lastErr
	}
	return ProviderResponse{}, fmt.Errorf("openai request failed after %d attempts over %s: %w", attempts, time.Since(startedAt).Round(time.Second), lastErr)
}

func (provider OpenAIProvider) debugf(format string, args ...interface{}) {
	if provider.DebugWriter == nil {
		return
	}
	if _, err := fmt.Fprintf(provider.DebugWriter, "[openai] "+format+"\n", args...); err != nil {
		return
	}
}

// generateOnce performs one Responses API attempt. Retry policy stays in
// Generate so a single attempt has straightforward HTTP and JSON ownership.
//
// Intent: Keep Flex retry decisions explicit and bounded without mixing retry
// state into response parsing. Source: DI-mopob
func (provider OpenAIProvider) generateOnce(ctx context.Context, request ProviderRequest, attempt int) (ProviderResponse, error) {
	if provider.APIKey == "" {
		return ProviderResponse{}, fmt.Errorf("missing OpenAI API key")
	}
	serviceTier, err := normalizeServiceTier(request.ServiceTier)
	if err != nil {
		return ProviderResponse{}, err
	}
	baseURL := provider.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/responses"
	}
	client := provider.Client
	if client == nil {
		// Intent: Do not allow direct OpenAIProvider construction to fall back to
		// the old 30-minute request wait. Source: DI-juzus
		client = &http.Client{Timeout: defaultRequestTimeout}
	}
	body := openAIRequest{
		Model:           request.APIModel,
		Input:           request.Prompt,
		Instructions:    request.Instructions,
		ServiceTier:     serviceTier,
		MaxOutputTokens: request.MaxOutputTokens,
	}
	if request.ReasoningEffort != "" {
		body.Reasoning = &openAIReasoning{Effort: request.ReasoningEffort}
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return ProviderResponse{}, err
	}
	// Intent: Make provider stalls diagnosable from the terminal canary log by
	// emitting the exact Responses API JSON sent for each attempt, without
	// logging authorization headers. Source: DI-juzus
	provider.debugf("attempt=%d event=request method=%s url=%s query_json=%s", attempt, http.MethodPost, baseURL, string(payload))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(payload))
	if err != nil {
		return ProviderResponse{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+provider.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	startedAt := time.Now()
	httpResp, err := client.Do(httpReq)
	if err != nil {
		provider.debugf("attempt=%d event=http_error elapsed=%s error=%q", attempt, time.Since(startedAt).Round(time.Millisecond), err.Error())
		return ProviderResponse{}, err
	}
	responseBytes, err := io.ReadAll(httpResp.Body)
	closeErr := httpResp.Body.Close()
	if err != nil {
		provider.debugf("attempt=%d event=response_read_error elapsed=%s status=%d request_id=%q error=%q", attempt, time.Since(startedAt).Round(time.Millisecond), httpResp.StatusCode, httpResp.Header.Get("x-request-id"), err.Error())
		return ProviderResponse{}, err
	}
	if closeErr != nil {
		provider.debugf("attempt=%d event=response_close_error elapsed=%s status=%d request_id=%q error=%q", attempt, time.Since(startedAt).Round(time.Millisecond), httpResp.StatusCode, httpResp.Header.Get("x-request-id"), closeErr.Error())
		return ProviderResponse{}, closeErr
	}
	// Intent: Preserve raw provider response evidence in the canary transcript
	// so empty-output, queued, incomplete, or tier-mismatch responses can be
	// diagnosed after an unattended run. Source: DI-juzus
	provider.debugf("attempt=%d event=response elapsed=%s status=%d request_id=%q response_json=%s", attempt, time.Since(startedAt).Round(time.Millisecond), httpResp.StatusCode, httpResp.Header.Get("x-request-id"), strings.TrimSpace(string(responseBytes)))
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return ProviderResponse{}, openAIHTTPError{StatusCode: httpResp.StatusCode, Body: strings.TrimSpace(string(responseBytes))}
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
		// Intent: Empty successful Responses payloads are transient provider
		// anomalies in canary evidence; retry before the caller skips the cell.
		// Source: DI-zikag
		return ProviderResponse{}, openAIResponseError{Message: "openai response contained no output text", Retryable: true}
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
		// Intent: Incomplete Responses payloads may succeed on retry and should
		// not abort an unattended GA cycle until bounded retries are exhausted.
		// Source: DI-zikag
		return ProviderResponse{}, openAIResponseError{
			Message:   fmt.Sprintf("openai response status %q with output text length %d", parsed.Status, len(text)),
			Retryable: parsed.Status == "incomplete",
		}
	}
	return ProviderResponse{
		Text:        strings.TrimSpace(text) + "\n",
		RequestID:   httpResp.Header.Get("x-request-id"),
		ResponseID:  parsed.ID,
		ServiceTier: parsed.ServiceTier,
		UsageJSON:   usageJSON,
	}, nil
}

func (policy ProviderRetryPolicy) withDefaults() ProviderRetryPolicy {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = defaultProviderMaxAttempts
	}
	if policy.MaxElapsed <= 0 {
		policy.MaxElapsed = defaultProviderMaxElapsed
	}
	if policy.InitialBackoff <= 0 {
		policy.InitialBackoff = 15 * time.Second
	}
	if policy.MaxBackoff <= 0 {
		policy.MaxBackoff = 2 * time.Minute
	}
	return policy
}

func (policy ProviderRetryPolicy) backoff(attempt int) time.Duration {
	backoff := policy.InitialBackoff
	for counter := 1; counter < attempt; counter++ {
		if backoff >= policy.MaxBackoff/2 {
			return policy.MaxBackoff
		}
		backoff *= 2
	}
	if backoff > policy.MaxBackoff {
		return policy.MaxBackoff
	}
	return backoff
}

type openAIHTTPError struct {
	StatusCode int
	Body       string
}

func (err openAIHTTPError) Error() string {
	if err.Body == "" {
		return fmt.Sprintf("openai status %d", err.StatusCode)
	}
	return fmt.Sprintf("openai status %d: %s", err.StatusCode, err.Body)
}

type openAIResponseError struct {
	Message   string
	Retryable bool
}

func (err openAIResponseError) Error() string {
	return err.Message
}

func shouldRetryOpenAIError(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	var responseErr openAIResponseError
	if errors.As(err, &responseErr) {
		return responseErr.Retryable
	}
	var httpErr openAIHTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusTooManyRequests || httpErr.StatusCode == http.StatusRequestTimeout
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return errors.Is(err, context.DeadlineExceeded)
}

func sleepWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

type openAIRequest struct {
	Model           string           `json:"model"`
	Input           string           `json:"input"`
	Instructions    string           `json:"instructions,omitempty"`
	ServiceTier     string           `json:"service_tier"`
	MaxOutputTokens int              `json:"max_output_tokens,omitempty"`
	Reasoning       *openAIReasoning `json:"reasoning,omitempty"`
}

type openAIReasoning struct {
	Effort string `json:"effort"`
}

type openAIResponse struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	OutputText  string                 `json:"output_text"`
	ServiceTier string                 `json:"service_tier"`
	Output      []openAIOutputItem     `json:"output"`
	Usage       map[string]interface{} `json:"usage"`
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
