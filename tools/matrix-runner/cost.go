package main

import (
	"encoding/json"
	"fmt"
)

const (
	defaultInputUSDPerMTok       = 1.75
	defaultCachedInputUSDPerMTok = 0.175
	defaultOutputUSDPerMTok      = 14.00
)

// CostConfig carries provider pricing and budget limits for an unattended run.
//
// Intent: Make matrix-run cost explicit and enforceable before expensive full
// runs start, while keeping the pricing values as CLI data rather than hard
// coding a provider pricing page snapshot into result semantics. Source:
// DI-nugiv
type CostConfig struct {
	InputUSDPerMTok       float64
	CachedInputUSDPerMTok float64
	OutputUSDPerMTok      float64
	MaxRunUSD             float64
	MaxCellEstimateUSD    float64
}

// UsageCost is the normalized token and dollar accounting recorded in queue
// state after a provider returns usage metadata.
type UsageCost struct {
	InputTokens       int
	CachedInputTokens int
	OutputTokens      int
	CostUSD           float64
}

type providerUsage struct {
	InputTokens        int                  `json:"input_tokens"`
	InputTokensDetails providerInputDetails `json:"input_tokens_details"`
	OutputTokens       int                  `json:"output_tokens"`
}

type providerInputDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

func defaultCostConfig() CostConfig {
	return CostConfig{
		InputUSDPerMTok:       defaultInputUSDPerMTok,
		CachedInputUSDPerMTok: defaultCachedInputUSDPerMTok,
		OutputUSDPerMTok:      defaultOutputUSDPerMTok,
	}
}

// BudgetEnabled reports whether missing usage metadata should be treated as a
// hard error because a run budget depends on accurate provider accounting.
func (config CostConfig) BudgetEnabled() bool {
	return config.MaxRunUSD > 0 || config.MaxCellEstimateUSD > 0
}

// CostForTokens applies the configured price table to normalized token counts.
// Cached input is clamped to total input so malformed provider metadata cannot
// produce a negative uncached-input total.
func (config CostConfig) CostForTokens(inputTokens int, cachedInputTokens int, outputTokens int) float64 {
	if cachedInputTokens > inputTokens {
		cachedInputTokens = inputTokens
	}
	uncachedInputTokens := inputTokens - cachedInputTokens
	return (float64(uncachedInputTokens)*config.InputUSDPerMTok +
		float64(cachedInputTokens)*config.CachedInputUSDPerMTok +
		float64(outputTokens)*config.OutputUSDPerMTok) / 1_000_000
}

// ParseUsage converts provider-specific usage JSON into the runner's normalized
// token/cost shape. The current OpenAI Responses usage schema is the first
// supported format.
func (config CostConfig) ParseUsage(usageJSON string) (UsageCost, error) {
	if usageJSON == "" {
		return UsageCost{}, fmt.Errorf("missing provider usage metadata")
	}
	var usage providerUsage
	if err := json.Unmarshal([]byte(usageJSON), &usage); err != nil {
		return UsageCost{}, err
	}
	cachedTokens := usage.InputTokensDetails.CachedTokens
	return UsageCost{
		InputTokens:       usage.InputTokens,
		CachedInputTokens: cachedTokens,
		OutputTokens:      usage.OutputTokens,
		CostUSD:           config.CostForTokens(usage.InputTokens, cachedTokens, usage.OutputTokens),
	}, nil
}

// EstimatePromptCost gives a conservative preflight budget check: it treats all
// prompt input as uncached and assumes the response may consume the configured
// output cap.
func (config CostConfig) EstimatePromptCost(prompt string, maxOutputTokens int) UsageCost {
	inputTokens := estimateTokenCount(prompt)
	return UsageCost{
		InputTokens:  inputTokens,
		OutputTokens: maxOutputTokens,
		CostUSD:      config.CostForTokens(inputTokens, 0, maxOutputTokens),
	}
}

// estimateTokenCount intentionally uses a rough character heuristic rather than
// a provider tokenizer. It is only a preflight guardrail; actual cost comes from
// provider usage metadata after a cell completes.
func estimateTokenCount(text string) int {
	runes := len([]rune(text))
	if runes == 0 {
		return 0
	}
	return (runes + 3) / 4
}

// stateActualCostUSD totals already-recorded provider costs from queue state so
// restarted runs keep respecting the original budget.
func stateActualCostUSD(state *QueueState) float64 {
	total := 0.0
	defaultPrices := defaultCostConfig()
	for _, record := range state.Cells {
		if record.CostUSD > 0 {
			total += record.CostUSD
			continue
		}
		if record.UsageJSON == "" {
			continue
		}
		if usageCost, err := defaultPrices.ParseUsage(record.UsageJSON); err == nil {
			total += usageCost.CostUSD
		}
	}
	return total
}
