package main

import (
	"context"
	"fmt"
	"strings"
)

const (
	defaultServiceTier = "flex"
	serviceTierFlex    = "flex"
	serviceTierDefault = "default"
)

// ProviderRequest is the complete prompt envelope that a GA command sends to a
// model provider for one score cell or one child-generation step.
type ProviderRequest struct {
	Provider        string
	APIModel        string
	ReasoningEffort string
	ServiceTier     string
	MaxOutputTokens int
	Instructions    string
	Prompt          string
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
