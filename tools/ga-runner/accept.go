package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type acceptOptions struct {
	RepoRoot     string
	RunGroupID   string
	ChildIDs     []string
	ResultPaths  []string
	ReviewerNote string
}

type stringListFlag []string

func (flagValue *stringListFlag) String() string {
	return strings.Join(*flagValue, ",")
}

func (flagValue *stringListFlag) Set(value string) error {
	if value == "" {
		return fmt.Errorf("empty value is not allowed")
	}
	*flagValue = append(*flagValue, value)
	return nil
}

type acceptancePlan struct {
	StatePath    string
	Children     map[string]GAChild
	ResultPaths  []string
	StagePaths   []string
	AcceptedAt   string
	ReviewerNote string
}

func runAccept(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("accept", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	runGroupID := fs.String("run-group-id", "", "GA run group whose state file records the child")
	reviewerNote := fs.String("reviewer-note", "", "human review note explaining why selected children are accepted")
	var childIDs stringListFlag
	var resultPaths stringListFlag
	fs.Var(&childIDs, "child", "child simulation ID to accept; repeatable")
	fs.Var(&resultPaths, "result", "selected JSON fitness result path; repeatable")
	if err := fs.Parse(args); err != nil {
		return err
	}
	options := acceptOptions{
		RepoRoot:     *repoRoot,
		RunGroupID:   *runGroupID,
		ChildIDs:     []string(childIDs),
		ResultPaths:  []string(resultPaths),
		ReviewerNote: *reviewerNote,
	}
	repo, err := openRepo(options.RepoRoot)
	if err != nil {
		return err
	}
	stateFile, err := statePath(repo, options.RunGroupID)
	if err != nil {
		return err
	}
	state, err := readGAState(stateFile)
	if err != nil {
		return err
	}
	plan, err := validateAcceptance(repo, stateFile, state, options)
	if err != nil {
		return err
	}
	now := plan.AcceptedAt
	state.UpdatedAt = now
	state.Acceptance = append(state.Acceptance, AcceptanceRecord{
		AcceptedChildIDs: sortedKeys(plan.Children),
		ResultPaths:      plan.ResultPaths,
		ReviewerNote:     plan.ReviewerNote,
		AcceptedAt:       now,
	})
	for index := range state.Children {
		childID := state.Children[index].ID()
		if _, ok := plan.Children[childID]; ok {
			state.Children[index].Status = "accepted"
		}
	}
	if err := writeGAStateAtomic(stateFile, state); err != nil {
		return err
	}
	return printAcceptance(stdout, plan)
}

// validateAcceptance verifies that review selects real generated children and
// JSON fitness evidence before the state file records acceptance.
//
// Intent: Promotion must be driven by reviewed state and result evidence, not by
// whatever untracked directories happen to exist under simulations. Source:
// DI-podot
func validateAcceptance(repo Repo, stateFile string, state GAState, options acceptOptions) (acceptancePlan, error) {
	if state.RunGroupID != options.RunGroupID {
		return acceptancePlan{}, fmt.Errorf("state run_group_id %q does not match requested %q", state.RunGroupID, options.RunGroupID)
	}
	if len(options.ChildIDs) == 0 {
		return acceptancePlan{}, fmt.Errorf("at least one -child is required")
	}
	if len(options.ResultPaths) == 0 {
		return acceptancePlan{}, fmt.Errorf("at least one -result is required")
	}
	if strings.TrimSpace(options.ReviewerNote) == "" {
		return acceptancePlan{}, fmt.Errorf("reviewer-note is required")
	}
	childrenByID, err := selectedChildren(repo, state, options.ChildIDs)
	if err != nil {
		return acceptancePlan{}, err
	}
	resultPaths, err := selectedResults(repo, state, childrenByID, options.ResultPaths)
	if err != nil {
		return acceptancePlan{}, err
	}
	stagePaths := []string{repo.Rel(stateFile)}
	for _, childID := range sortedKeys(childrenByID) {
		stagePaths = append(stagePaths, normalizeRelPath(childrenByID[childID].Path))
	}
	stagePaths = append(stagePaths, resultPaths...)
	return acceptancePlan{
		StatePath:    repo.Rel(stateFile),
		Children:     childrenByID,
		ResultPaths:  resultPaths,
		StagePaths:   uniqueStrings(stagePaths),
		AcceptedAt:   time.Now().UTC().Format(time.RFC3339),
		ReviewerNote: strings.TrimSpace(options.ReviewerNote),
	}, nil
}

func selectedChildren(repo Repo, state GAState, childIDs []string) (map[string]GAChild, error) {
	stateChildren := map[string]GAChild{}
	for _, child := range state.Children {
		childID := child.ID()
		if childID != "" {
			stateChildren[childID] = child
		}
	}
	selected := map[string]GAChild{}
	for _, childID := range childIDs {
		if _, exists := selected[childID]; exists {
			continue
		}
		child, ok := stateChildren[childID]
		if !ok {
			return nil, fmt.Errorf("child %s is not present in GA state", childID)
		}
		if err := validateChildForAcceptance(repo, state, childID, child); err != nil {
			return nil, err
		}
		selected[childID] = child
	}
	return selected, nil
}

func validateChildForAcceptance(repo Repo, state GAState, childID string, child GAChild) error {
	if !strings.HasPrefix(childID, "SIM-") {
		return fmt.Errorf("child %s must start with SIM-", childID)
	}
	if child.Status == "culled" {
		return fmt.Errorf("child %s is already culled", childID)
	}
	childPath := strings.TrimSuffix(normalizeRelPath(child.Path), "/")
	expectedPath := strings.TrimSuffix(proposalChildSimulationPath(state.RunGroupID, childID), "/")
	if childPath != expectedPath {
		return fmt.Errorf("child %s path must be %s", childID, expectedPath+"/")
	}
	info, err := os.Stat(repo.Abs(childPath))
	if err != nil {
		return fmt.Errorf("child %s path is not readable: %w", childID, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("child %s path is not a directory", childID)
	}
	currentHash, err := currentSimulationTreeHash(repo, childPath)
	if err != nil {
		return err
	}
	if child.TreeHash == "" {
		return fmt.Errorf("child %s state tree_hash is required", childID)
	}
	if currentHash != child.TreeHash {
		return fmt.Errorf("child %s tree hash drift: state=%s current=%s", childID, child.TreeHash, currentHash)
	}
	return nil
}

func selectedResults(repo Repo, state GAState, children map[string]GAChild, resultPaths []string) ([]string, error) {
	// Intent: Keep acceptance tied to JSON fitness artifacts from this GA state
	// rather than old Markdown canaries or unrelated result files. Source: DI-podot
	stateResultPaths := stateCellResultPaths(state)
	var selected []string
	for _, inputPath := range resultPaths {
		absPath := repo.Abs(inputPath)
		relPath := normalizeRelPath(repo.Rel(absPath))
		parts, pathIssues := parseResultPath(repo, absPath)
		if len(pathIssues) > 0 {
			return nil, fmt.Errorf("%s: %s", relPath, strings.Join(pathIssues, "; "))
		}
		if _, ok := children[parts.SimID]; !ok {
			return nil, fmt.Errorf("result %s belongs to %s, which is not a selected child", relPath, parts.SimID)
		}
		if len(state.Cells) > 0 && !stateResultPaths[relPath] {
			return nil, fmt.Errorf("result %s is not present in GA state cells", relPath)
		}
		if issues := validateResultFile(repo, absPath); len(issues) > 0 {
			return nil, fmt.Errorf("%s: %s", relPath, strings.Join(issues, "; "))
		}
		result, err := readFitnessResult(absPath)
		if err != nil {
			return nil, err
		}
		child := children[parts.SimID]
		if err := validateResultForChild(repo, relPath, result, child); err != nil {
			return nil, err
		}
		selected = append(selected, relPath)
	}
	return uniqueStrings(selected), nil
}

func validateResultForChild(repo Repo, resultPath string, result FitnessResult, child GAChild) error {
	childPath := strings.TrimSuffix(normalizeRelPath(child.Path), "/")
	sourcePath := strings.TrimSuffix(normalizeRelPath(result.Source.SimPath), "/")
	if sourcePath != childPath {
		return fmt.Errorf("result %s source.sim_path %s does not match child path %s", resultPath, sourcePath, childPath)
	}
	currentHash, err := currentSimulationTreeHash(repo, childPath)
	if err != nil {
		return err
	}
	if result.Source.SimulationTreeHash != currentHash {
		return fmt.Errorf("result %s source simulation tree hash does not match current child tree", resultPath)
	}
	return nil
}

func stateCellResultPaths(state GAState) map[string]bool {
	paths := map[string]bool{}
	for _, cell := range state.Cells {
		if path := cell.SelectedResultPath(); path != "" {
			paths[normalizeRelPath(path)] = true
		}
	}
	return paths
}

func currentSimulationTreeHash(repo Repo, relDir string) (string, error) {
	// Intent: Re-hash the current materialized child tree before acceptance so
	// reviewers do not promote a tree that drifted after scoring; proposal-stage
	// child paths are valid until promotion moves winners to canonical
	// `simulations/`. Source: DI-podot; DI-lirat
	cleanDir := strings.TrimSuffix(normalizeRelPath(relDir), "/")
	if !isCanonicalSimulationTreePath(cleanDir) && !isProposalSimulationTreePath(cleanDir) {
		return "", fmt.Errorf("simulation tree must be under simulations/SIM-* or proposals/<run-group>/simulations/SIM-*")
	}
	var files []string
	err := filepath.WalkDir(repo.Abs(cleanDir), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		files = append(files, repo.Rel(path))
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(files)
	if len(files) == 0 {
		return "", fmt.Errorf("simulation tree %s has no files", cleanDir)
	}
	return trackedTreeHash(repo, files)
}

func isCanonicalSimulationTreePath(relPath string) bool {
	parts := strings.Split(strings.TrimSuffix(normalizeRelPath(relPath), "/"), "/")
	return len(parts) == 2 && parts[0] == "simulations" && safeSimIDPattern.MatchString(parts[1])
}

func isProposalSimulationTreePath(relPath string) bool {
	parts := strings.Split(strings.TrimSuffix(normalizeRelPath(relPath), "/"), "/")
	return len(parts) == 4 &&
		parts[0] == "proposals" &&
		safeStateIDPattern.MatchString(parts[1]) &&
		parts[2] == "simulations" &&
		safeSimIDPattern.MatchString(parts[3])
}

func readFitnessResult(path string) (FitnessResult, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return FitnessResult{}, err
	}
	var result FitnessResult
	if err := json.Unmarshal(bytes, &result); err != nil {
		return FitnessResult{}, err
	}
	return result, nil
}

func printAcceptance(stdout io.Writer, plan acceptancePlan) error {
	// Intent: Report exact staging paths while leaving git index and commit
	// control to the operator's normal repository workflow. Source: DI-podot
	if err := writeFormat(stdout, "accepted children=%s results=%d state=%s\n", strings.Join(sortedKeys(plan.Children), ","), len(plan.ResultPaths), plan.StatePath); err != nil {
		return err
	}
	if err := writeLine(stdout, "review paths:"); err != nil {
		return err
	}
	for _, path := range plan.StagePaths {
		if err := writeFormat(stdout, "  %s\n", path); err != nil {
			return err
		}
	}
	// Intent: Proposal paths are ignored staging artifacts; acceptance records
	// reviewer intent, while a later promotion step must move selected artifacts
	// into canonical `simulations/` and `results/` paths before git staging.
	// Source: DI-lirat
	return writeLine(stdout, "promotion required before git add: move selected proposal artifacts into canonical simulations/ and results/ paths")
}

func sortedKeys(values map[string]GAChild) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	var unique []string
	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		unique = append(unique, value)
	}
	return unique
}
