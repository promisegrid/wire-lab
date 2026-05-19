package main

import "context"

type ProviderRequest struct {
	Provider        string
	APIModel        string
	ReasoningEffort string
	MaxOutputTokens int
	Prompt          string
}

type ProviderResponse struct {
	Text       string
	RequestID  string
	ResponseID string
	UsageJSON  string
}

type Provider interface {
	Generate(ctx context.Context, request ProviderRequest) (ProviderResponse, error)
}
