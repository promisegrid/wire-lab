package main

import "context"

// ProviderRequest is the complete prompt envelope that a GA command sends to a
// model provider for one score cell or one child-generation step.
type ProviderRequest struct {
	Provider        string
	APIModel        string
	ReasoningEffort string
	MaxOutputTokens int
	Instructions    string
	Prompt          string
}

// ProviderResponse is intentionally narrow: the runner consumes the model text
// plus provider IDs and raw usage metadata, then fills authoritative GA result
// fields itself.
type ProviderResponse struct {
	Text       string
	RequestID  string
	ResponseID string
	UsageJSON  string
}

// Provider hides provider-specific APIs behind the single GA operation needed
// in v1: generate structured text for score or child-bundle prompts.
type Provider interface {
	Generate(ctx context.Context, request ProviderRequest) (ProviderResponse, error)
}
