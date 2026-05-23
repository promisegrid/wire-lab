package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const stateSchemaV1 = "promisegrid.ga.state.v1"

var safeStateIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// GAState is the v1 checkpoint file that makes pending children auditable.
//
// Intent: Keep generated child sims outside the stable population until review
// records an explicit acceptance event in the run state. Source: DI-podot
type GAState struct {
	Schema         string             `json:"schema"`
	RunGroupID     string             `json:"run_group_id"`
	CreatedAt      string             `json:"created_at,omitempty"`
	UpdatedAt      string             `json:"updated_at,omitempty"`
	RepoCommit     string             `json:"repo_commit,omitempty"`
	ModelID        string             `json:"model_id,omitempty"`
	Population     []GAStateSim       `json:"population,omitempty"`
	ScenarioSample []GAStateScenario  `json:"scenario_sample,omitempty"`
	Parents        []GAStateParent    `json:"parents,omitempty"`
	Children       []GAChild          `json:"children,omitempty"`
	Cells          []GACell           `json:"cells,omitempty"`
	Acceptance     []AcceptanceRecord `json:"acceptance,omitempty"`
	Culling        []CullingRecord    `json:"culling,omitempty"`
}

type GAStateSim struct {
	SimID    string `json:"sim_id"`
	Path     string `json:"path"`
	TreeHash string `json:"tree_hash"`
}

type GAStateScenario struct {
	ScenarioID   string `json:"scenario_id"`
	Path         string `json:"path"`
	SamplePolicy string `json:"sample_policy,omitempty"`
	SHA256       string `json:"sha256,omitempty"`
	TreeHash     string `json:"tree_hash,omitempty"`
}

type GAStateParent struct {
	SimID     string `json:"sim_id"`
	Rationale string `json:"rationale,omitempty"`
}

// GAChild describes a generated child simulation proposal under proposals/.
//
// Intent: Acceptance verifies the durable child tree instead of trusting a JSON
// proposal object or the mere presence of a directory. Source: DI-podot
type GAChild struct {
	ChildID            string       `json:"child_id,omitempty"`
	SimID              string       `json:"sim_id,omitempty"`
	Path               string       `json:"path"`
	ParentIDs          []string     `json:"parent_ids,omitempty"`
	Operation          string       `json:"operation,omitempty"`
	PromptHash         string       `json:"prompt_hash,omitempty"`
	ResponseHash       string       `json:"response_hash,omitempty"`
	Files              []SourceFile `json:"files,omitempty"`
	TreeHash           string       `json:"tree_hash"`
	Status             string       `json:"status"`
	DesignDeltaSummary string       `json:"design_delta_summary,omitempty"`
	ValidationMessage  string       `json:"validation_message,omitempty"`
	UpdatedAt          string       `json:"updated_at,omitempty"`
	Provider           string       `json:"provider,omitempty"`
	APIModel           string       `json:"api_model,omitempty"`
	ReasoningEffort    string       `json:"reasoning_effort,omitempty"`
	ServiceTier        string       `json:"service_tier,omitempty"`
	ServedServiceTier  string       `json:"served_service_tier,omitempty"`
	RequestID          string       `json:"request_id,omitempty"`
	ResponseID         string       `json:"response_id,omitempty"`
	UsageJSON          string       `json:"usage_json,omitempty"`
	InputTokens        int          `json:"input_tokens,omitempty"`
	CachedTokens       int          `json:"cached_input_tokens,omitempty"`
	OutputTokens       int          `json:"output_tokens,omitempty"`
	CostUSD            float64      `json:"cost_usd,omitempty"`
}

type GACell struct {
	CellID             string  `json:"cell_id"`
	SimID              string  `json:"sim_id"`
	ScenarioID         string  `json:"scenario_id"`
	ModelID            string  `json:"model_id"`
	ResultPath         string  `json:"result_path,omitempty"`
	ExpectedResultPath string  `json:"expected_result_path,omitempty"`
	Status             string  `json:"status,omitempty"`
	Attempts           int     `json:"attempts,omitempty"`
	ValidationMessage  string  `json:"validation_message,omitempty"`
	UpdatedAt          string  `json:"updated_at,omitempty"`
	Provider           string  `json:"provider,omitempty"`
	APIModel           string  `json:"api_model,omitempty"`
	ReasoningEffort    string  `json:"reasoning_effort,omitempty"`
	ServiceTier        string  `json:"service_tier,omitempty"`
	ServedServiceTier  string  `json:"served_service_tier,omitempty"`
	RequestID          string  `json:"request_id,omitempty"`
	ResponseID         string  `json:"response_id,omitempty"`
	UsageJSON          string  `json:"usage_json,omitempty"`
	InputTokens        int     `json:"input_tokens,omitempty"`
	CachedTokens       int     `json:"cached_input_tokens,omitempty"`
	OutputTokens       int     `json:"output_tokens,omitempty"`
	CostUSD            float64 `json:"cost_usd,omitempty"`
}

type AcceptanceRecord struct {
	AcceptedChildIDs []string `json:"accepted_child_ids"`
	ResultPaths      []string `json:"result_paths"`
	ReviewerNote     string   `json:"reviewer_note"`
	AcceptedAt       string   `json:"accepted_at"`
}

type CullingRecord struct {
	CulledChildIDs     []string `json:"culled_child_ids,omitempty"`
	DeletedSimPaths    []string `json:"deleted_sim_paths,omitempty"`
	DeletedResultPaths []string `json:"deleted_result_paths,omitempty"`
	Reason             string   `json:"reason,omitempty"`
	CulledAt           string   `json:"culled_at,omitempty"`
}

func statePath(repo Repo, runGroupID string) (string, error) {
	// Intent: Treat run group IDs as filenames only, preventing a review command
	// from resolving state outside results/state/. Source: DI-podot
	if runGroupID == "" {
		return "", fmt.Errorf("run-group-id is required")
	}
	if !safeStateIDPattern.MatchString(runGroupID) {
		return "", fmt.Errorf("run-group-id must be a safe path segment")
	}
	return repo.Path("results", "state", runGroupID+".json"), nil
}

func readGAState(path string) (GAState, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return GAState{}, err
	}
	var state GAState
	if err := json.Unmarshal(bytes, &state); err != nil {
		return GAState{}, err
	}
	if state.Schema != stateSchemaV1 {
		return GAState{}, fmt.Errorf("state schema must be %s", stateSchemaV1)
	}
	return state, nil
}

// writeGAStateAtomic writes a v1 GA state update through a temp file in the same
// directory so interrupted review checkpoints do not leave truncated JSON.
//
// Intent: Acceptance is a review checkpoint, so the state file must move from
// old to new content atomically. Source: DI-podot
func writeGAStateAtomic(path string, state GAState) error {
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	if err := ensureParent(path); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(bytes, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (child GAChild) ID() string {
	if child.ChildID != "" {
		return child.ChildID
	}
	return child.SimID
}

func (cell GACell) SelectedResultPath() string {
	if cell.ResultPath != "" {
		return filepath.ToSlash(cell.ResultPath)
	}
	return filepath.ToSlash(cell.ExpectedResultPath)
}

func normalizeRelPath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

// deriveAPIModelFromModelID peels provider and reasoning suffix decorations off
// the durable result-path model ID so mixed historical backfill runs can reuse
// the provider model each cell originally came from. Source: DI-roruj
func deriveAPIModelFromModelID(modelID string) string {
	clean := strings.TrimSpace(strings.ToLower(modelID))
	if clean == "" {
		return ""
	}
	clean = strings.TrimPrefix(clean, "openai-")
	for _, suffix := range []string{"-xhigh", "-high", "-medium", "-low"} {
		if strings.HasSuffix(clean, suffix) {
			return strings.TrimSuffix(clean, suffix)
		}
	}
	return clean
}

func proposalRunRoot(runGroupID string) string {
	return filepath.ToSlash(filepath.Join("proposals", runGroupID))
}

func proposalChildSimulationPath(runGroupID string, childID string) string {
	// Intent: Keep generated child simulation trees out of canonical
	// `simulations/` until human review promotes them. Source: DI-lirat
	return filepath.ToSlash(filepath.Join(proposalRunRoot(runGroupID), "simulations", childID)) + "/"
}

func proposalChildResultRoot(runGroupID string, childID string) string {
	return filepath.ToSlash(filepath.Join(proposalRunRoot(runGroupID), "results", childID))
}

func proposalChildResultPath(runGroupID string, childID string, scenarioID string, modelID string, timestamp string) string {
	// Intent: Keep generated child score evidence with the ignored child proposal
	// it evaluates, rather than mixing it into canonical `results/`. Source:
	// DI-lirat
	return filepath.ToSlash(filepath.Join(proposalChildResultRoot(runGroupID, childID), scenarioID, modelID, timestamp+".json"))
}

func isRunScopedProposalSimulationPath(runGroupID string, childID string, relPath string) bool {
	expected := strings.TrimSuffix(normalizeRelPath(proposalChildSimulationPath(runGroupID, childID)), "/")
	return strings.TrimSuffix(normalizeRelPath(relPath), "/") == expected
}

func simulationPathForState(state GAState, simID string) string {
	for _, child := range state.Children {
		if child.ID() == simID && child.Path != "" {
			// Intent: Child score cells must read the ignored proposal tree recorded
			// in GA state, not an accidental canonical tree with the same ID.
			// Source: DI-lirat
			return strings.TrimSuffix(normalizeRelPath(child.Path), "/")
		}
	}
	return filepath.ToSlash(filepath.Join("simulations", simID))
}
