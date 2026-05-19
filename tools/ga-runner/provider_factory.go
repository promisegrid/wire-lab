package main

import (
	"fmt"
	"os"
)

// buildProvider centralizes provider construction so score/generate share the
// same dry-run and API-key behavior.
//
// Intent: V1 GA runs support one OpenAI-compatible provider path while keeping
// the command code provider-neutral for later extension. Source: DI-gijom
func buildProvider(providerName string, apiKeyEnv string, openAIBaseURL string, dryRun bool) (Provider, error) {
	if dryRun {
		return nil, nil
	}
	switch providerName {
	case "openai":
		return OpenAIProvider{
			APIKey:  os.Getenv(apiKeyEnv),
			BaseURL: openAIBaseURL,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", providerName)
	}
}
