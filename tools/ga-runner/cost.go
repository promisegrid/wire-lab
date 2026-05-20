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

const (
	// Intent: Keep preflight budget estimates conservative after removing
	// provider hard output caps from default score/generate requests. Source:
	// DI-pulap
	defaultScoreCostEstimateOutputTokens    = 4000
	defaultGenerateCostEstimateOutputTokens = 6000
)

// CostConfig keeps expensive GA runs bounded by explicit operator budgets.
//
// Intent: Provider-backed GA scoring and generation can become expensive, so
// score/generate need the same preflight and measured-cost controls as earlier
// matrix runs. Source: DI-gijom
type CostConfig struct {
	InputUSDPerMTok       float64
	CachedInputUSDPerMTok float64
	OutputUSDPerMTok      float64
	MaxRunUSD             float64
	MaxCellEstimateUSD    float64
}

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

func (config CostConfig) BudgetEnabled() bool {
	return config.MaxRunUSD > 0 || config.MaxCellEstimateUSD > 0
}

func (config CostConfig) CostForTokens(inputTokens int, cachedInputTokens int, outputTokens int) float64 {
	if cachedInputTokens > inputTokens {
		cachedInputTokens = inputTokens
	}
	uncachedInputTokens := inputTokens - cachedInputTokens
	return (float64(uncachedInputTokens)*config.InputUSDPerMTok +
		float64(cachedInputTokens)*config.CachedInputUSDPerMTok +
		float64(outputTokens)*config.OutputUSDPerMTok) / 1_000_000
}

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

func (config CostConfig) EstimatePromptCost(prompt string, maxOutputTokens int) UsageCost {
	inputTokens := estimateTokenCount(prompt)
	return UsageCost{
		InputTokens:  inputTokens,
		OutputTokens: maxOutputTokens,
		CostUSD:      config.CostForTokens(inputTokens, 0, maxOutputTokens),
	}
}

func estimateTokenCount(text string) int {
	runes := len([]rune(text))
	if runes == 0 {
		return 0
	}
	return (runes + 3) / 4
}

func stateActualCostUSD(state GAState) float64 {
	total := 0.0
	defaultPrices := CostConfig{
		InputUSDPerMTok:       defaultInputUSDPerMTok,
		CachedInputUSDPerMTok: defaultCachedInputUSDPerMTok,
		OutputUSDPerMTok:      defaultOutputUSDPerMTok,
	}
	for _, cell := range state.Cells {
		if cell.CostUSD > 0 {
			total += cell.CostUSD
			continue
		}
		if cell.UsageJSON == "" {
			continue
		}
		if usageCost, err := defaultPrices.ParseUsage(cell.UsageJSON); err == nil {
			total += usageCost.CostUSD
		}
	}
	for _, child := range state.Children {
		if child.CostUSD > 0 {
			total += child.CostUSD
			continue
		}
		if child.UsageJSON == "" {
			continue
		}
		if usageCost, err := defaultPrices.ParseUsage(child.UsageJSON); err == nil {
			total += usageCost.CostUSD
		}
	}
	return total
}
