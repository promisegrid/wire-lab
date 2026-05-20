package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type providerBuildOptions struct {
	ProviderName        string
	APIKeyEnv           string
	OpenAIBaseURL       string
	DryRun              bool
	RequestTimeout      time.Duration
	ProviderMaxAttempts int
	ProviderMaxElapsed  time.Duration
	Stream              bool
	StreamIdleTimeout   time.Duration
}

// buildProvider centralizes provider construction so score/generate share the
// same dry-run and API-key behavior.
//
// Intent: V1 GA runs support one OpenAI-compatible provider path while keeping
// the command code provider-neutral for later extension. Source: DI-gijom
func buildProvider(options providerBuildOptions) (Provider, error) {
	if options.DryRun {
		return nil, nil
	}
	requestTimeout := options.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = defaultRequestTimeout
	}
	maxAttempts := options.ProviderMaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = defaultProviderMaxAttempts
	}
	maxElapsed := options.ProviderMaxElapsed
	if maxElapsed <= 0 {
		maxElapsed = defaultProviderMaxElapsed
	}
	streamIdleTimeout := options.StreamIdleTimeout
	if streamIdleTimeout <= 0 {
		streamIdleTimeout = defaultStreamIdleTimeout
	}
	switch options.ProviderName {
	case "openai":
		// Intent: Bound each synchronous API attempt inside OpenAIProvider rather
		// than with http.Client.Timeout so streaming events can reset the idle
		// watchdog while the retry loop still has a real second attempt budget.
		// Source: DI-tufud
		return OpenAIProvider{
			APIKey:            os.Getenv(options.APIKeyEnv),
			BaseURL:           options.OpenAIBaseURL,
			Client:            &http.Client{},
			RequestTimeout:    requestTimeout,
			Stream:            options.Stream,
			StreamIdleTimeout: streamIdleTimeout,
			// Intent: Send raw provider request/response diagnostics to the
			// canary transcript so hangs have evidence before timeout. Source:
			// DI-tufud
			DebugWriter: os.Stderr,
			RetryPolicy: ProviderRetryPolicy{
				MaxAttempts: maxAttempts,
				MaxElapsed:  maxElapsed,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", options.ProviderName)
	}
}
