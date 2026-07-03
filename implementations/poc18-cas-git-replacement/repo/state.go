package repo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

const defaultCurrentBranch = "main"

// State is the `.grid/state.json` shape for local mutable CLI conveniences.
//
// Intent: Keep the current local branch and snapshot pointers outside
// `.grid/config.json` so configuration remains stable while commands can still
// default to the user's current local graph position. Source: DI-bikif
type State struct {
	Version                   int      `json:"version"`
	CurrentBranch             string   `json:"current_branch"`
	CurrentSnapshotCID        string   `json:"current_snapshot_cid,omitempty"`
	CurrentBranchRefSetCID    string   `json:"current_branch_refset_cid,omitempty"`
	CurrentWorkspaceRefSetCID string   `json:"current_workspace_refset_cid,omitempty"`
	UntrackedPaths            []string `json:"untracked_paths,omitempty"`
}

func defaultState() State {
	return State{Version: 1, CurrentBranch: defaultCurrentBranch}
}

// LoadState reads `.grid/state.json`, returning a default empty state when the
// file has not yet been created.
func (repository Repository) LoadState() (State, error) {
	content, readErr := os.ReadFile(repository.statePath())
	if os.IsNotExist(readErr) {
		return defaultState(), nil
	}
	if readErr != nil {
		return State{}, readErr
	}
	return parseState(content)
}

// SaveState validates and writes `.grid/state.json` atomically.
func (repository Repository) SaveState(state State) error {
	if state.Version == 0 {
		state.Version = 1
	}
	if state.CurrentBranch == "" {
		state.CurrentBranch = defaultCurrentBranch
	}
	normalizedPaths, pathErr := normalizeStatePaths(state.UntrackedPaths)
	if pathErr != nil {
		return pathErr
	}
	state.UntrackedPaths = normalizedPaths
	if validateErr := validateState(state); validateErr != nil {
		return validateErr
	}
	content, marshalErr := json.MarshalIndent(state, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	content = append(content, '\n')
	if mkdirErr := os.MkdirAll(repository.GridDir, 0o755); mkdirErr != nil {
		return mkdirErr
	}
	statePath := repository.statePath()
	tmpPath := statePath + ".tmp"
	if writeErr := os.WriteFile(tmpPath, content, 0o644); writeErr != nil {
		return writeErr
	}
	return os.Rename(tmpPath, statePath)
}

// RecordSnapshot updates local mutable state after a successful snapshot.
func (repository Repository) RecordSnapshot(snapshotCID, branchRefSetCID, workspaceRefSetCID string) (State, error) {
	state, loadErr := repository.LoadState()
	if loadErr != nil {
		return State{}, loadErr
	}
	state.CurrentSnapshotCID = snapshotCID
	state.CurrentBranchRefSetCID = branchRefSetCID
	state.CurrentWorkspaceRefSetCID = workspaceRefSetCID
	if saveErr := repository.SaveState(state); saveErr != nil {
		return State{}, saveErr
	}
	return state, nil
}

// TrackPaths removes repo-relative paths from local exclusions.
//
// Intent: `grid track` is not Git staging. It simply lets paths re-enter the
// default POC18 policy where every non-control workspace path is considered by
// snapshot and status. Source: DI-jokav
func (repository Repository) TrackPaths(paths []string) (State, error) {
	state, loadErr := repository.LoadState()
	if loadErr != nil {
		return State{}, loadErr
	}
	excluded := statePathSet(state.UntrackedPaths)
	for _, path := range paths {
		normalized, pathErr := NormalizeStatePath(path)
		if pathErr != nil {
			return State{}, pathErr
		}
		delete(excluded, normalized)
	}
	state.UntrackedPaths = sortedStatePaths(excluded)
	if saveErr := repository.SaveState(state); saveErr != nil {
		return State{}, saveErr
	}
	return state, nil
}

// UntrackPaths records repo-relative paths as local exclusions.
//
// Intent: `grid untrack` omits local-only or generated paths from future
// repo-local snapshots and status reports without asserting that peers should
// share the same policy. Source: DI-jokav
func (repository Repository) UntrackPaths(paths []string) (State, error) {
	state, loadErr := repository.LoadState()
	if loadErr != nil {
		return State{}, loadErr
	}
	excluded := statePathSet(state.UntrackedPaths)
	for _, path := range paths {
		normalized, pathErr := NormalizeStatePath(path)
		if pathErr != nil {
			return State{}, pathErr
		}
		excluded[normalized] = true
	}
	state.UntrackedPaths = sortedStatePaths(excluded)
	if saveErr := repository.SaveState(state); saveErr != nil {
		return State{}, saveErr
	}
	return state, nil
}

func parseState(content []byte) (State, error) {
	var state State
	if unmarshalErr := json.Unmarshal(content, &state); unmarshalErr != nil {
		return State{}, unmarshalErr
	}
	normalizedPaths, pathErr := normalizeStatePaths(state.UntrackedPaths)
	if pathErr != nil {
		return State{}, pathErr
	}
	state.UntrackedPaths = normalizedPaths
	if validateErr := validateState(state); validateErr != nil {
		return State{}, validateErr
	}
	return state, nil
}

func validateState(state State) error {
	if state.Version != 1 {
		return fmt.Errorf("unsupported state version %d", state.Version)
	}
	if state.CurrentBranch == "" {
		return fmt.Errorf("current_branch is required")
	}
	for field, cidText := range map[string]string{
		"current_snapshot_cid":         state.CurrentSnapshotCID,
		"current_branch_refset_cid":    state.CurrentBranchRefSetCID,
		"current_workspace_refset_cid": state.CurrentWorkspaceRefSetCID,
	} {
		if cidText == "" {
			continue
		}
		if _, parseErr := store.ParseCIDText(cidText); parseErr != nil {
			return fmt.Errorf("%s: %w", field, parseErr)
		}
	}
	if _, pathErr := normalizeStatePaths(state.UntrackedPaths); pathErr != nil {
		return pathErr
	}
	return nil
}

func (repository Repository) statePath() string {
	return filepath.Join(repository.GridDir, StateFileName)
}

// NormalizeStatePath returns the canonical printable form for a path stored in
// `.grid/state.json`.
//
// Intent: State paths are local repo-relative policy keys, not global resource
// identifiers. Keeping them slash-clean avoids platform-specific drift while
// rejecting attempts to make `.grid` control state part of versioned content.
// Source: DI-jokav
func NormalizeStatePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("path must be repo-relative: %s", path)
	}
	slashPath := filepath.ToSlash(path)
	for _, segment := range strings.Split(slashPath, "/") {
		if segment == ".." {
			return "", fmt.Errorf("path must not contain parent traversal: %s", path)
		}
	}
	cleaned := filepath.ToSlash(filepath.Clean(path))
	if cleaned == "" || cleaned == "." {
		return "", fmt.Errorf("path must name a workspace entry")
	}
	if cleaned == GridDirName || strings.HasPrefix(cleaned, GridDirName+"/") {
		return "", fmt.Errorf("path must not target %s control state: %s", GridDirName, path)
	}
	return cleaned, nil
}

func normalizeStatePaths(paths []string) ([]string, error) {
	seen := map[string]bool{}
	for _, path := range paths {
		normalized, pathErr := NormalizeStatePath(path)
		if pathErr != nil {
			return nil, pathErr
		}
		seen[normalized] = true
	}
	return sortedStatePaths(seen), nil
}

func statePathSet(paths []string) map[string]bool {
	seen := map[string]bool{}
	for _, path := range paths {
		seen[path] = true
	}
	return seen
}

func sortedStatePaths(paths map[string]bool) []string {
	sorted := make([]string, 0, len(paths))
	for path := range paths {
		sorted = append(sorted, path)
	}
	sort.Strings(sorted)
	return sorted
}
