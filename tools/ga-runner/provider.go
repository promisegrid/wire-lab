package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	defaultServiceTier = "flex"
	serviceTierFlex    = "flex"
	serviceTierDefault = "default"
)

const (
	defaultRequestTimeout      = 5 * time.Minute
	defaultProviderMaxAttempts = 2
	// Intent: Leave enough elapsed retry budget for two full five-minute
	// attempts plus backoff instead of timing out the second attempt early.
	// Source: DI-tufud
	defaultProviderMaxElapsed = 12 * time.Minute
	// Intent: Prefer streaming liveness diagnostics for canary provider calls
	// unless the operator explicitly disables streaming. Source: DI-tufud
	defaultProviderStream = true
	// Intent: Retry a silent streaming connection while allowing active SSE
	// events to prove the provider request is still alive. Source: DI-tufud
	defaultStreamIdleTimeout = 2 * time.Minute
	defaultScoreWorkers      = 1
	defaultGenerateWorkers   = 1
)

const (
	// Intent: Score operations keep the established high-reasoning comparison
	// baseline, while child generation defaults lower because canary evidence
	// showed xhigh child responses exhausting output caps on hidden reasoning.
	// Source: DI-pulap
	defaultScoreReasoningEffort    = "xhigh"
	defaultGenerateReasoningEffort = "medium"
	defaultTextVerbosity           = "low"
	textVerbosityLow               = "low"
	textVerbosityMedium            = "medium"
	textVerbosityHigh              = "high"
)

// ProviderRequest is the complete prompt envelope that a GA command sends to a
// model provider for one score cell or one child-generation step.
type ProviderRequest struct {
	Provider         string
	APIModel         string
	ReasoningEffort  string
	ReasoningSummary string
	ServiceTier      string
	TextVerbosity    string
	MaxOutputTokens  int
	Instructions     string
	Prompt           string
}

// ProviderResponse is intentionally narrow: the runner consumes the model text
// plus provider IDs and raw usage metadata, then fills authoritative GA result
// fields itself.
type ProviderResponse struct {
	Text        string
	RequestID   string
	ResponseID  string
	ServiceTier string
	UsageJSON   string
}

// Provider hides provider-specific APIs behind the single GA operation needed
// in v1: generate structured text for score or child-bundle prompts.
type Provider interface {
	Generate(ctx context.Context, request ProviderRequest) (ProviderResponse, error)
}

// normalizeServiceTier keeps provider-backed GA runs on an explicit cost tier.
//
// Intent: Prevent expensive Priority or project-inherited processing from
// entering unattended GA runs by allowing only explicit Flex or Standard
// requests. Source: DI-mopob
func normalizeServiceTier(serviceTier string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(serviceTier))
	if normalized == "" {
		return defaultServiceTier, nil
	}
	switch normalized {
	case serviceTierFlex, serviceTierDefault:
		return normalized, nil
	default:
		return "", fmt.Errorf("service-tier must be %q or %q; got %q", serviceTierFlex, serviceTierDefault, serviceTier)
	}
}

// normalizeTextVerbosity keeps concise JSON shaping explicit while avoiding a
// hard output-token cap that can exhaust on hidden reasoning tokens.
//
// Intent: Guide result size with the provider's soft text-verbosity control
// instead of default `max_output_tokens` caps that made child generation fail
// after spending reasoning tokens. Source: DI-pulap
func normalizeTextVerbosity(textVerbosity string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(textVerbosity))
	if normalized == "" {
		return defaultTextVerbosity, nil
	}
	switch normalized {
	case textVerbosityLow, textVerbosityMedium, textVerbosityHigh:
		return normalized, nil
	default:
		return "", fmt.Errorf("text-verbosity must be %q, %q, or %q; got %q", textVerbosityLow, textVerbosityMedium, textVerbosityHigh, textVerbosity)
	}
}

func normalizeWorkers(workers int) (int, error) {
	if workers < 1 {
		return 0, fmt.Errorf("workers must be at least 1")
	}
	return workers, nil
}

func parsePositiveDurationFlag(name string, value string) (time.Duration, error) {
	duration, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("%s must be a Go duration such as 5m: %w", name, err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", name)
	}
	return duration, nil
}
