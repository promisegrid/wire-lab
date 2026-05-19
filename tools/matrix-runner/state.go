package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type QueueState struct {
	Manifest   string                `json:"manifest"`
	RunGroupID string                `json:"run_group_id"`
	CreatedAt  string                `json:"created_at"`
	UpdatedAt  string                `json:"updated_at"`
	Cells      map[string]*CellState `json:"cells"`
}

type CellState struct {
	CellID          string `json:"cell_id"`
	Ordinal         int    `json:"ordinal"`
	SimID           string `json:"sim_id"`
	ScenarioID      string `json:"scenario_id"`
	ModelID         string `json:"model_id"`
	ResultPath      string `json:"result_path"`
	Status          string `json:"status"`
	Attempts        int    `json:"attempts"`
	LastMessage     string `json:"last_message"`
	UpdatedAt       string `json:"updated_at"`
	Provider        string `json:"provider,omitempty"`
	APIModel        string `json:"api_model,omitempty"`
	ReasoningEffort string `json:"reasoning_effort,omitempty"`
	RequestID       string `json:"request_id,omitempty"`
	ResponseID      string `json:"response_id,omitempty"`
	// Intent: Keep measured provider usage with each cell so budget restarts,
	// progress checks, and run audits do not depend on re-reading result prose.
	// Source: DI-nugiv
	UsageJSON    string  `json:"usage_json,omitempty"`
	InputTokens  int     `json:"input_tokens,omitempty"`
	CachedTokens int     `json:"cached_input_tokens,omitempty"`
	OutputTokens int     `json:"output_tokens,omitempty"`
	CostUSD      float64 `json:"cost_usd,omitempty"`
}

func utcISO() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func defaultStatePath(repo Repo, cells []MatrixCell) string {
	return repo.Path("results", "state", slug(cells[0].RunGroupID)+".json")
}

func defaultJobDir(repo Repo, cells []MatrixCell) string {
	return repo.Path("results", "jobs", cells[0].RunGroupID)
}

func loadOrCreateState(repo Repo, statePath string, manifestPath string, cells []MatrixCell) (*QueueState, error) {
	if bytes, err := os.ReadFile(statePath); err == nil {
		var state QueueState
		if err := json.Unmarshal(bytes, &state); err != nil {
			return nil, err
		}
		if state.Cells == nil {
			state.Cells = map[string]*CellState{}
		}
		mergeStateCells(&state, cells)
		return &state, nil
	}
	state := &QueueState{
		Manifest:   repo.Rel(manifestPath),
		RunGroupID: cells[0].RunGroupID,
		CreatedAt:  utcISO(),
		Cells:      map[string]*CellState{},
	}
	mergeStateCells(state, cells)
	return state, nil
}

func mergeStateCells(state *QueueState, cells []MatrixCell) {
	for _, cell := range cells {
		record := state.Cells[cell.CellID]
		if record == nil {
			record = &CellState{
				CellID:      cell.CellID,
				Status:      cell.Status,
				LastMessage: "",
			}
			if record.Status == "" {
				record.Status = "queued"
			}
			state.Cells[cell.CellID] = record
		}
		record.Ordinal = cell.Ordinal
		record.SimID = cell.SimID
		record.ScenarioID = cell.ScenarioID
		record.ModelID = cell.ModelID
		record.ResultPath = cell.ResultPath
	}
}

func saveState(path string, state *QueueState) error {
	state.UpdatedAt = utcISO()
	if err := ensureParent(path); err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(bytes, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func markCell(record *CellState, status string, message string) {
	record.Status = status
	record.LastMessage = message
	record.UpdatedAt = utcISO()
}

func stateCounts(state *QueueState) map[string]int {
	counts := map[string]int{}
	for _, record := range state.Cells {
		status := record.Status
		if status == "" {
			status = "queued"
		}
		counts[status]++
	}
	return counts
}

func formatCounts(state *QueueState) string {
	counts := stateCounts(state)
	var statuses []string
	for status := range counts {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)
	total := 0
	for _, count := range counts {
		total += count
	}
	parts := []string{"total=" + strconv.Itoa(total)}
	for _, status := range statuses {
		parts = append(parts, status+"="+strconv.Itoa(counts[status]))
	}
	return strings.Join(parts, " ")
}

func cleanPath(path string) string {
	return filepath.Clean(path)
}
