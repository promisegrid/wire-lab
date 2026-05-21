package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var safeSimIDPattern = regexp.MustCompile(`^SIM-[A-Za-z0-9._-]+$`)

type cullOptions struct {
	RepoRoot   string
	RunGroupID string
	ChildIDs   []string
	Reason     string
	DryRun     bool
}

type cullingPlan struct {
	StatePath          string
	Children           map[string]GAChild
	DeletedSimPaths    []string
	DeletedResultPaths []string
	MissingResultPaths []string
	Reason             string
	CulledAt           string
	DryRun             bool
}

func runCull(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("cull", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	runGroupID := fs.String("run-group-id", "", "GA run group whose state file records generated children")
	reason := fs.String("reason", "", "human review note explaining why selected children are rejected")
	dryRun := fs.Bool("dry-run", false, "validate and print culling targets without deleting or writing state")
	var childIDs stringListFlag
	fs.Var(&childIDs, "child", "child simulation ID to cull; repeatable")
	if err := fs.Parse(args); err != nil {
		return err
	}
	options := cullOptions{
		RepoRoot:   *repoRoot,
		RunGroupID: *runGroupID,
		ChildIDs:   []string(childIDs),
		Reason:     *reason,
		DryRun:     *dryRun,
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
	plan, err := validateCulling(repo, stateFile, state, options)
	if err != nil {
		return err
	}
	if options.DryRun {
		return printCullingPlan(stdout, plan)
	}
	if err := deleteCullingTargets(repo, plan); err != nil {
		return err
	}
	now := plan.CulledAt
	state.UpdatedAt = now
	state.Culling = append(state.Culling, CullingRecord{
		CulledChildIDs:     sortedKeys(plan.Children),
		DeletedSimPaths:    plan.DeletedSimPaths,
		DeletedResultPaths: plan.DeletedResultPaths,
		Reason:             plan.Reason,
		CulledAt:           now,
	})
	for index := range state.Children {
		childID := state.Children[index].ID()
		if _, ok := plan.Children[childID]; ok {
			state.Children[index].Status = "culled"
		}
	}
	if err := writeGAStateAtomic(stateFile, state); err != nil {
		return err
	}
	return printCullingPlan(stdout, plan)
}

// validateCulling builds a deletion plan only from generated children recorded
// in the v1 GA state file.
//
// Intent: Culling is intentionally destructive, so target discovery must be
// state-bound and exact instead of filesystem-wide. Source: DI-kofil
func validateCulling(repo Repo, stateFile string, state GAState, options cullOptions) (cullingPlan, error) {
	if state.RunGroupID != options.RunGroupID {
		return cullingPlan{}, fmt.Errorf("state run_group_id %q does not match requested %q", state.RunGroupID, options.RunGroupID)
	}
	if len(options.ChildIDs) == 0 {
		return cullingPlan{}, fmt.Errorf("at least one -child is required")
	}
	if strings.TrimSpace(options.Reason) == "" {
		return cullingPlan{}, fmt.Errorf("reason is required")
	}
	children, err := selectedChildrenForCulling(repo, state, options.ChildIDs)
	if err != nil {
		return cullingPlan{}, err
	}
	var simPaths []string
	var resultPaths []string
	var missingResultPaths []string
	for _, childID := range sortedKeys(children) {
		childPath := strings.TrimSuffix(normalizeRelPath(children[childID].Path), "/")
		resultPath := proposalChildResultRoot(state.RunGroupID, childID)
		simPaths = append(simPaths, childPath)
		if info, err := os.Stat(repo.Abs(resultPath)); err == nil && info.IsDir() {
			resultPaths = append(resultPaths, resultPath)
		} else if err == nil {
			return cullingPlan{}, fmt.Errorf("%s exists but is not a directory", resultPath)
		} else if os.IsNotExist(err) {
			missingResultPaths = append(missingResultPaths, resultPath)
		} else {
			return cullingPlan{}, err
		}
	}
	return cullingPlan{
		StatePath:          repo.Rel(stateFile),
		Children:           children,
		DeletedSimPaths:    uniqueStrings(simPaths),
		DeletedResultPaths: uniqueStrings(resultPaths),
		MissingResultPaths: uniqueStrings(missingResultPaths),
		Reason:             strings.TrimSpace(options.Reason),
		CulledAt:           time.Now().UTC().Format(time.RFC3339),
		DryRun:             options.DryRun,
	}, nil
}

func selectedChildrenForCulling(repo Repo, state GAState, childIDs []string) (map[string]GAChild, error) {
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
		if err := validateChildForCulling(repo, state, childID, child); err != nil {
			return nil, err
		}
		selected[childID] = child
	}
	return selected, nil
}

func validateChildForCulling(repo Repo, state GAState, childID string, child GAChild) error {
	if !safeSimIDPattern.MatchString(childID) {
		return fmt.Errorf("child %s must be a safe SIM-* path segment", childID)
	}
	if child.Status == "accepted" {
		return fmt.Errorf("child %s is accepted and cannot be culled", childID)
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
	return nil
}

// deleteCullingTargets removes only the state-validated child sim and matching
// result directories from an already-built culling plan.
//
// Intent: Keep destructive cleanup narrow and replayable from the culling record
// rather than deriving targets during deletion. Source: DI-kofil
func deleteCullingTargets(repo Repo, plan cullingPlan) error {
	targets := append([]string{}, plan.DeletedSimPaths...)
	targets = append(targets, plan.DeletedResultPaths...)
	for _, relPath := range targets {
		if err := validateDeletionTarget(relPath); err != nil {
			return err
		}
		if err := os.RemoveAll(repo.Abs(relPath)); err != nil {
			return err
		}
	}
	return nil
}

func validateDeletionTarget(relPath string) error {
	clean := normalizeRelPath(relPath)
	if filepath.IsAbs(relPath) || strings.HasPrefix(clean, "..") {
		return fmt.Errorf("refusing unsafe deletion target %s", relPath)
	}
	if isCanonicalSimulationTreePath(clean) {
		return nil
	}
	parts := strings.Split(clean, "/")
	if len(parts) == 4 &&
		parts[0] == "proposals" &&
		safeStateIDPattern.MatchString(parts[1]) &&
		(parts[2] == "simulations" || parts[2] == "results") &&
		safeSimIDPattern.MatchString(parts[3]) {
		return nil
	}
	return fmt.Errorf("refusing deletion target outside proposal child sim/result roots: %s", relPath)
}

func printCullingPlan(stdout io.Writer, plan cullingPlan) error {
	// Intent: Make destructive cleanup visible before and after execution,
	// including dry-run plans that make no filesystem changes. Source: DI-kofil
	mode := "cull"
	if plan.DryRun {
		mode = "dry-run"
	}
	if err := writeFormat(stdout, "%s children=%s state=%s\n", mode, strings.Join(sortedKeys(plan.Children), ","), plan.StatePath); err != nil {
		return err
	}
	if err := writeLine(stdout, "delete simulation paths:"); err != nil {
		return err
	}
	for _, path := range sortedCopy(plan.DeletedSimPaths) {
		if err := writeFormat(stdout, "  %s\n", path); err != nil {
			return err
		}
	}
	if err := writeLine(stdout, "delete result paths:"); err != nil {
		return err
	}
	for _, path := range sortedCopy(plan.DeletedResultPaths) {
		if err := writeFormat(stdout, "  %s\n", path); err != nil {
			return err
		}
	}
	for _, path := range sortedCopy(plan.MissingResultPaths) {
		if err := writeFormat(stdout, "  missing %s\n", path); err != nil {
			return err
		}
	}
	return nil
}

func sortedCopy(values []string) []string {
	copied := append([]string{}, values...)
	sort.Strings(copied)
	return copied
}
