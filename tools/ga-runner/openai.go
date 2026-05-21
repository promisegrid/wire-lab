package main

import (
	"bufio"
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
	APIKey  string
	BaseURL string
	Client  *http.Client
	// Intent: Bound one provider attempt while preserving a larger total retry
	// window for a useful second attempt. Source: DI-tufud
	RequestTimeout time.Duration
	Stream         bool
	// Intent: Treat silent streaming gaps as retryable stalls while any SSE event
	// proves the provider connection is still alive. Source: DI-tufud
	StreamIdleTimeout time.Duration
	// Intent: Mirror selected streaming content to stdout for canary diagnosis
	// without changing the parsed provider response. Source: DI-vadub
	StreamContentWriter io.Writer
	RetryPolicy         ProviderRetryPolicy
	// DebugWriter receives raw request and response diagnostics for provider
	// calls. Authorization headers are intentionally excluded from these logs.
	DebugWriter io.Writer
}

// ProviderRetryPolicy bounds provider retries for long-running unattended GA
// runs. The zero value selects conservative sync canary defaults.
//
// Intent: Keep retry timing explicit, finite, and testable so Flex scarcity and
// slow provider calls do not require interactive babysitting. Source: DI-mopob;
// DI-tufud
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
// checkpointable provider error and streams liveness events when enabled.
// Source: DI-gijom; DI-mopob; DI-tufud
func (provider OpenAIProvider) Generate(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
	policy := provider.RetryPolicy.withDefaults()
	var lastErr error
	startedAt := time.Now()
	attempts := 0
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		attempts = attempt
		attemptCtx := ctx
		cancelAttempt := func() {}
		if timeout := provider.requestTimeout(); timeout > 0 {
			// Intent: Scope the request timeout to this attempt, not to the whole
			// retry window, so a timed-out first attempt still leaves the second
			// attempt with its own full timeout budget. Source: DI-tufud
			attemptCtx, cancelAttempt = context.WithTimeout(ctx, timeout)
		}
		response, err := provider.generateOnce(attemptCtx, request, attempt)
		cancelAttempt()
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
	textVerbosity, err := normalizeTextVerbosity(request.TextVerbosity)
	if err != nil {
		return ProviderResponse{}, err
	}
	baseURL := provider.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/responses"
	}
	client := provider.Client
	if client == nil {
		// Intent: Keep timeout ownership in the attempt context so streaming can
		// use idle diagnostics without http.Client cutting off active SSE
		// connections independently. Source: DI-tufud
		client = &http.Client{}
	}
	body := openAIRequest{
		Model:           request.APIModel,
		Input:           request.Prompt,
		Instructions:    request.Instructions,
		ServiceTier:     serviceTier,
		MaxOutputTokens: request.MaxOutputTokens,
		Stream:          provider.Stream,
		Text:            &openAIText{Verbosity: textVerbosity},
	}
	if request.ReasoningEffort != "" || request.ReasoningSummary != "" {
		body.Reasoning = &openAIReasoning{Effort: request.ReasoningEffort, Summary: request.ReasoningSummary}
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
	if provider.Stream {
		// Intent: Use Responses API streaming events as liveness evidence while
		// preserving the final JSON parsing contract for score and child bundle
		// responses. Source: DI-tufud
		return provider.readStreamingResponse(ctx, httpResp, attempt, startedAt)
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
	return normalizeOpenAIResponse(parsed, httpResp.Header.Get("x-request-id"), "")
}

func normalizeOpenAIResponse(parsed openAIResponse, requestID string, textOverride string) (ProviderResponse, error) {
	usageJSON := ""
	if len(parsed.Usage) > 0 {
		usageBytes, err := json.Marshal(parsed.Usage)
		if err != nil {
			return ProviderResponse{}, err
		}
		usageJSON = string(usageBytes)
	}
	text := parsed.OutputText
	if textOverride != "" {
		text = textOverride
	} else if text == "" {
		text = parsed.JoinOutputText()
	}
	if parsed.Error != nil {
		// Intent: Preserve provider-side failed Response evidence while treating
		// ordinary server-side failures as retryable provider anomalies.
		// Source: DI-tufud
		return ProviderResponse{}, openAIResponseError{
			Message:   fmt.Sprintf("openai response error type %q code %q message %q usage %s", parsed.Error.Type, parsed.Error.Code, parsed.Error.Message, usageJSON),
			Retryable: retryableOpenAIProviderError(parsed.Error.Type),
		}
	}
	if parsed.Status != "" && parsed.Status != "completed" {
		reason := parsed.IncompleteDetails.Reason
		// Intent: Retry transient incomplete Responses states, but do not retry
		// deterministic max-output exhaustion. The canary log showed xhigh
		// calls consuming the whole output cap as reasoning tokens; retrying the
		// same cap only spends more budget and looks like a hang. Source:
		// DI-juzus; DI-zikag
		return ProviderResponse{}, openAIResponseError{
			Message: fmt.Sprintf("openai response status %q reason %q with output text length %d usage %s",
				parsed.Status, reason, len(text), usageJSON),
			Retryable: retryableOpenAIStatus(parsed.Status, reason),
		}
	}
	if strings.TrimSpace(text) == "" {
		// Intent: Empty successful Responses payloads are transient provider
		// anomalies in canary evidence; retry before the caller skips the cell.
		// Source: DI-zikag
		return ProviderResponse{}, openAIResponseError{
			Message:   fmt.Sprintf("openai response contained no output text usage %s", usageJSON),
			Retryable: true,
		}
	}
	return ProviderResponse{
		Text:        strings.TrimSpace(text) + "\n",
		RequestID:   requestID,
		ResponseID:  parsed.ID,
		ServiceTier: parsed.ServiceTier,
		UsageJSON:   usageJSON,
	}, nil
}

func (provider OpenAIProvider) readStreamingResponse(ctx context.Context, httpResp *http.Response, attempt int, startedAt time.Time) (ProviderResponse, error) {
	defer func() {
		if err := httpResp.Body.Close(); err != nil {
			provider.debugf("attempt=%d event=response_close_error elapsed=%s status=%d request_id=%q error=%q", attempt, time.Since(startedAt).Round(time.Millisecond), httpResp.StatusCode, httpResp.Header.Get("x-request-id"), err.Error())
		}
	}()
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		responseBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			provider.debugf("attempt=%d event=response_read_error elapsed=%s status=%d request_id=%q error=%q", attempt, time.Since(startedAt).Round(time.Millisecond), httpResp.StatusCode, httpResp.Header.Get("x-request-id"), err.Error())
			return ProviderResponse{}, err
		}
		// Intent: Preserve raw non-2xx provider response evidence for retry and
		// postmortem analysis even when streaming was requested. Source:
		// DI-tufud
		provider.debugf("attempt=%d event=response elapsed=%s status=%d request_id=%q response_json=%s", attempt, time.Since(startedAt).Round(time.Millisecond), httpResp.StatusCode, httpResp.Header.Get("x-request-id"), strings.TrimSpace(string(responseBytes)))
		return ProviderResponse{}, openAIHTTPError{StatusCode: httpResp.StatusCode, Body: strings.TrimSpace(string(responseBytes))}
	}
	idleTimeout := provider.streamIdleTimeout()
	lines := scanOpenAIStreamLines(httpResp.Body)
	timer := time.NewTimer(idleTimeout)
	defer timer.Stop()
	output := &strings.Builder{}
	eventType := ""
	var dataLines []string
	for {
		select {
		case <-ctx.Done():
			return ProviderResponse{}, ctx.Err()
		case <-timer.C:
			// Intent: Turn a silent streaming connection into retryable evidence
			// instead of leaving the canary blocked with no observable progress.
			// Source: DI-tufud
			return ProviderResponse{}, openAIResponseError{
				Message:   fmt.Sprintf("openai stream idle for %s after %s", idleTimeout, time.Since(startedAt).Round(time.Millisecond)),
				Retryable: true,
			}
		case lineResult, ok := <-lines:
			if !ok {
				if len(dataLines) > 0 {
					response, done, err := provider.dispatchOpenAIStreamEvent(attempt, startedAt, eventType, strings.Join(dataLines, "\n"), httpResp.Header.Get("x-request-id"), output)
					if done || err != nil {
						return response, err
					}
				}
				return ProviderResponse{}, openAIResponseError{
					Message:   fmt.Sprintf("openai stream ended before response.completed after %s", time.Since(startedAt).Round(time.Millisecond)),
					Retryable: true,
				}
			}
			resetTimer(timer, idleTimeout)
			if lineResult.Err != nil {
				return ProviderResponse{}, openAIResponseError{
					Message:   fmt.Sprintf("openai stream read error after %s: %v", time.Since(startedAt).Round(time.Millisecond), lineResult.Err),
					Retryable: true,
				}
			}
			line := strings.TrimSuffix(lineResult.Line, "\r")
			if line == "" {
				if len(dataLines) == 0 {
					eventType = ""
					continue
				}
				response, done, err := provider.dispatchOpenAIStreamEvent(attempt, startedAt, eventType, strings.Join(dataLines, "\n"), httpResp.Header.Get("x-request-id"), output)
				eventType = ""
				dataLines = nil
				if done || err != nil {
					return response, err
				}
				continue
			}
			if strings.HasPrefix(line, ":") {
				provider.debugf("attempt=%d event=stream_keepalive elapsed=%s bytes=%d", attempt, time.Since(startedAt).Round(time.Millisecond), len(line))
				continue
			}
			if after, ok := strings.CutPrefix(line, "event:"); ok {
				eventType = strings.TrimSpace(after)
				continue
			}
			if after, ok := strings.CutPrefix(line, "data:"); ok {
				dataLines = append(dataLines, strings.TrimSpace(after))
			}
		}
	}
}

type streamLineResult struct {
	Line string
	Err  error
}

func scanOpenAIStreamLines(reader io.Reader) <-chan streamLineResult {
	lines := make(chan streamLineResult, 64)
	go func() {
		defer close(lines)
		scanner := bufio.NewScanner(reader)
		scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
		for scanner.Scan() {
			lines <- streamLineResult{Line: scanner.Text()}
		}
		if err := scanner.Err(); err != nil {
			lines <- streamLineResult{Err: err}
		}
	}()
	return lines
}

func (provider OpenAIProvider) dispatchOpenAIStreamEvent(attempt int, startedAt time.Time, eventType string, data string, requestID string, output *strings.Builder) (ProviderResponse, bool, error) {
	if data == "[DONE]" {
		provider.debugf("attempt=%d event=stream_done_marker elapsed=%s", attempt, time.Since(startedAt).Round(time.Millisecond))
		return ProviderResponse{}, false, nil
	}
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal([]byte(data), &envelope); err != nil {
		return ProviderResponse{}, true, fmt.Errorf("decode openai stream event %q: %w", eventType, err)
	}
	if eventType == "" {
		eventType = envelope.Type
	}
	switch eventType {
	case "response.output_text.delta":
		var payload struct {
			Delta string `json:"delta"`
		}
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			return ProviderResponse{}, true, fmt.Errorf("decode openai stream text delta: %w", err)
		}
		output.WriteString(payload.Delta)
		return ProviderResponse{}, false, nil
	case "response.reasoning_summary_text.delta":
		var payload struct {
			Delta string `json:"delta"`
		}
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			return ProviderResponse{}, true, fmt.Errorf("decode openai stream reasoning summary delta: %w", err)
		}
		provider.writeStreamProgressDot(attempt, eventType)
		return ProviderResponse{}, false, nil
	case "response.reasoning_summary_part.done":
		var payload struct {
			Part struct {
				Text string `json:"text"`
			} `json:"part"`
		}
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			return ProviderResponse{}, true, fmt.Errorf("decode openai stream reasoning summary part: %w", err)
		}
		provider.writeStreamContent(attempt, eventType, payload.Part.Text)
		provider.debugf("attempt=%d event=stream_event elapsed=%s type=%q text_chars=%d", attempt, time.Since(startedAt).Round(time.Millisecond), eventType, len(payload.Part.Text))
		return ProviderResponse{}, false, nil
	case "response.output_text.done":
		var payload struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			return ProviderResponse{}, true, fmt.Errorf("decode openai stream text done: %w", err)
		}
		if output.Len() == 0 && payload.Text != "" {
			output.WriteString(payload.Text)
		}
		provider.debugf("attempt=%d event=stream_event elapsed=%s type=%q text_chars=%d", attempt, time.Since(startedAt).Round(time.Millisecond), eventType, len(payload.Text))
		return ProviderResponse{}, false, nil
	case "response.completed":
		var payload struct {
			Response openAIResponse `json:"response"`
		}
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			return ProviderResponse{}, true, fmt.Errorf("decode openai stream completed response: %w", err)
		}
		provider.debugf("attempt=%d event=stream_event elapsed=%s type=%q response_id=%q", attempt, time.Since(startedAt).Round(time.Millisecond), eventType, payload.Response.ID)
		response, err := normalizeOpenAIResponse(payload.Response, requestID, output.String())
		return response, true, err
	case "response.failed", "response.incomplete":
		var payload struct {
			Response openAIResponse `json:"response"`
		}
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			return ProviderResponse{}, true, fmt.Errorf("decode openai stream failed response: %w", err)
		}
		provider.debugf("attempt=%d event=stream_event elapsed=%s type=%q response_id=%q", attempt, time.Since(startedAt).Round(time.Millisecond), eventType, payload.Response.ID)
		response, err := normalizeOpenAIResponse(payload.Response, requestID, output.String())
		return response, true, err
	case "error":
		var payload struct {
			Error   *openAIError `json:"error"`
			Message string       `json:"message"`
			Type    string       `json:"type"`
			Code    string       `json:"code"`
		}
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			return ProviderResponse{}, true, fmt.Errorf("decode openai stream error: %w", err)
		}
		errorType := payload.Type
		errorCode := payload.Code
		message := payload.Message
		if payload.Error != nil {
			errorType = payload.Error.Type
			errorCode = payload.Error.Code
			message = payload.Error.Message
		}
		provider.debugf("attempt=%d event=stream_event elapsed=%s type=%q error_type=%q error_code=%q", attempt, time.Since(startedAt).Round(time.Millisecond), eventType, errorType, errorCode)
		return ProviderResponse{}, true, openAIResponseError{
			Message:   fmt.Sprintf("openai stream error type %q code %q message %q", errorType, errorCode, message),
			Retryable: retryableOpenAIProviderError(errorType),
		}
	default:
		provider.debugf("attempt=%d event=stream_event elapsed=%s type=%q data_bytes=%d", attempt, time.Since(startedAt).Round(time.Millisecond), eventType, len(data))
		return ProviderResponse{}, false, nil
	}
}

func resetTimer(timer *time.Timer, duration time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(duration)
}

func (provider OpenAIProvider) writeStreamContent(attempt int, eventType string, text string) {
	if provider.StreamContentWriter == nil || text == "" {
		return
	}
	// Intent: Write line-oriented completed reasoning-summary part events so
	// canary stdout/logs show useful summary checkpoints without corrupting the
	// provider response text used for JSON parsing. Source: DI-vadub; DI-babik;
	// DI-vajut; DI-sakam; DI-fupob; DI-ramun
	if _, err := fmt.Fprintf(provider.StreamContentWriter, "[openai-stream] attempt=%d type=%s delta=%q\n", attempt, eventType, text); err != nil {
		provider.debugf("attempt=%d event=stream_content_write_error type=%q error=%q", attempt, eventType, err.Error())
	}
}

func (provider OpenAIProvider) writeStreamProgressDot(attempt int, eventType string) {
	if provider.StreamContentWriter == nil {
		return
	}
	// Intent: Show that reasoning-summary streaming is alive without printing
	// summary event names or summary text into the canary transcript. Source:
	// DI-babik
	if _, err := fmt.Fprint(provider.StreamContentWriter, "."); err != nil {
		provider.debugf("attempt=%d event=stream_content_write_error type=%q error=%q", attempt, eventType, err.Error())
	}
}

func retryableOpenAIStatus(status string, reason string) bool {
	switch status {
	case "queued", "in_progress":
		return true
	case "incomplete":
		return reason != "max_output_tokens"
	default:
		return false
	}
}

func retryableOpenAIProviderError(errorType string) bool {
	switch strings.ToLower(strings.TrimSpace(errorType)) {
	case "invalid_request_error", "authentication_error", "permission_error", "insufficient_quota":
		return false
	default:
		return true
	}
}

func (provider OpenAIProvider) requestTimeout() time.Duration {
	if provider.RequestTimeout <= 0 {
		return defaultRequestTimeout
	}
	return provider.RequestTimeout
}

func (provider OpenAIProvider) streamIdleTimeout() time.Duration {
	if provider.StreamIdleTimeout <= 0 {
		return defaultStreamIdleTimeout
	}
	return provider.StreamIdleTimeout
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
		return retryableOpenAIHTTPStatus(httpErr.StatusCode)
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return errors.Is(err, context.DeadlineExceeded)
}

func retryableOpenAIHTTPStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
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
	Stream          bool             `json:"stream,omitempty"`
	Text            *openAIText      `json:"text,omitempty"`
	Reasoning       *openAIReasoning `json:"reasoning,omitempty"`
}

type openAIText struct {
	Verbosity string `json:"verbosity,omitempty"`
}

type openAIReasoning struct {
	Effort  string `json:"effort,omitempty"`
	Summary string `json:"summary,omitempty"`
}

type openAIResponse struct {
	ID                string                  `json:"id"`
	Status            string                  `json:"status"`
	OutputText        string                  `json:"output_text"`
	ServiceTier       string                  `json:"service_tier"`
	Error             *openAIError            `json:"error"`
	Output            []openAIOutputItem      `json:"output"`
	Usage             map[string]interface{}  `json:"usage"`
	IncompleteDetails openAIIncompleteDetails `json:"incomplete_details"`
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

type openAIIncompleteDetails struct {
	Reason string `json:"reason"`
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
