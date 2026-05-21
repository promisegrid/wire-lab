package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type SourceDocument struct {
	Path string
	Text string
}

// sha256File records source-document identity in GA result files so later
// reviews can detect sim/scenario drift.
func sha256File(repo Repo, relPath string) (string, error) {
	bytes, err := os.ReadFile(repo.Abs(relPath))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(bytes)
	return hex.EncodeToString(sum[:]), nil
}

func sourceFilesForResult(repo Repo, state GAState, simID string, scenario Scenario) ([]SourceFile, error) {
	paths, err := sourcePathsForPrompt(repo, simulationPathForState(state, simID), scenario)
	if err != nil {
		return nil, err
	}
	var files []SourceFile
	for _, path := range paths {
		hash, err := sha256File(repo, path)
		if err != nil {
			return nil, err
		}
		files = append(files, SourceFile{Path: path, SHA256: hash})
	}
	return files, nil
}

// sourcePathsForPrompt chooses the local contract, sim, and scenario documents
// that must be bundled into a provider prompt and hashed into the result.
//
// Intent: Keep each GA cell source-complete without restoring cross-sim shared
// source-of-truth files or relying on provider-side filesystem discovery.
// Proposal child cells use the state-recorded proposal path instead of canonical
// `simulations/`. Source: DI-gijom; DI-lirat
func sourcePathsForPrompt(repo Repo, simPath string, scenario Scenario) ([]string, error) {
	cleanSimPath := strings.TrimSuffix(normalizeRelPath(simPath), "/")
	required := []string{
		"results/RUN-PROTOCOL.md",
		"scenarios/README.md",
		filepath.ToSlash(filepath.Join(cleanSimPath, "README.md")),
	}
	var paths []string
	for _, rel := range required {
		if info, err := os.Stat(repo.Abs(rel)); err != nil || info.IsDir() {
			if err != nil {
				return nil, fmt.Errorf("read required source %s: %w", rel, err)
			}
			return nil, fmt.Errorf("required source %s is a directory", rel)
		}
		paths = append(paths, rel)
	}
	questionPath := filepath.ToSlash(filepath.Join(cleanSimPath, "QUESTION.md"))
	if info, err := os.Stat(repo.Abs(questionPath)); err == nil && !info.IsDir() {
		paths = append(paths, questionPath)
	}
	localSim, err := localMarkdownFiles(repo, cleanSimPath, map[string]bool{
		"README.md":   true,
		"QUESTION.md": true,
	})
	if err != nil {
		return nil, err
	}
	paths = append(paths, localSim...)
	paths = append(paths, scenario.Path)
	localScenario, err := localMarkdownFiles(repo, filepath.ToSlash(filepath.Join("scenarios", scenario.ScenarioID)), map[string]bool{
		filepath.Base(scenario.Path): true,
		"README.md":                  true,
		"MATRIX.md":                  true,
	})
	if err != nil {
		return nil, err
	}
	paths = append(paths, localScenario...)
	return uniqueStrings(paths), nil
}

func sourceDocumentsForPrompt(repo Repo, state GAState, simID string, scenario Scenario) ([]SourceDocument, error) {
	paths, err := sourcePathsForPrompt(repo, simulationPathForState(state, simID), scenario)
	if err != nil {
		return nil, err
	}
	return sourceDocumentsFromPaths(repo, paths)
}

func sourceDocumentsFromPaths(repo Repo, paths []string) ([]SourceDocument, error) {
	// Intent: Keep path-selection policy separate from file loading so scoring
	// and generation prompts can share safe repo-rooted reads. Source: DI-gijom;
	// DI-dilaf
	var docs []SourceDocument
	for _, path := range paths {
		text, err := repo.ReadRel(path)
		if err != nil {
			return nil, err
		}
		docs = append(docs, SourceDocument{Path: path, Text: text})
	}
	return docs, nil
}

func parentSourceDocumentsForGeneratePrompt(repo Repo, simID string) ([]SourceDocument, error) {
	// Intent: Keep child-generation prompts source-complete for each parent sim
	// while avoiding repeated root contracts and scenario boilerplate that made
	// provider calls time out before returning headers. Source: DI-dilaf
	required := []string{
		filepath.ToSlash(filepath.Join("simulations", simID, "README.md")),
	}
	var paths []string
	for _, rel := range required {
		if info, err := os.Stat(repo.Abs(rel)); err != nil || info.IsDir() {
			if err != nil {
				return nil, fmt.Errorf("read required source %s: %w", rel, err)
			}
			return nil, fmt.Errorf("required source %s is a directory", rel)
		}
		paths = append(paths, rel)
	}
	questionPath := filepath.ToSlash(filepath.Join("simulations", simID, "QUESTION.md"))
	if info, err := os.Stat(repo.Abs(questionPath)); err == nil && !info.IsDir() {
		paths = append(paths, questionPath)
	}
	localSim, err := localMarkdownFiles(repo, filepath.ToSlash(filepath.Join("simulations", simID)), map[string]bool{
		"README.md":   true,
		"QUESTION.md": true,
	})
	if err != nil {
		return nil, err
	}
	paths = append(paths, localSim...)
	return sourceDocumentsFromPaths(repo, uniqueStrings(paths))
}

func scenarioSourceDocumentsForGeneratePrompt(repo Repo, scenario GAStateScenario) ([]SourceDocument, error) {
	// Intent: Keep scenario-specific pressure complete for child generation while
	// excluding the shared root scenario contract that belongs in scoring prompts,
	// not every breed prompt. Source: DI-dilaf
	paths := []string{scenario.Path}
	localScenario, err := localMarkdownFiles(repo, filepath.ToSlash(filepath.Join("scenarios", scenario.ScenarioID)), map[string]bool{
		filepath.Base(scenario.Path): true,
		"README.md":                  true,
		"MATRIX.md":                  true,
	})
	if err != nil {
		return nil, err
	}
	paths = append(paths, localScenario...)
	return sourceDocumentsFromPaths(repo, uniqueStrings(paths))
}

func localMarkdownFiles(repo Repo, relRoot string, excludedBase map[string]bool) ([]string, error) {
	var paths []string
	root := repo.Abs(relRoot)
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s is not a directory", relRoot)
	}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		if excludedBase[filepath.Base(path)] {
			return nil
		}
		paths = append(paths, repo.Rel(path))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

func hashText(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

func pathSafeSlug(text string) string {
	lower := strings.ToLower(text)
	var out []rune
	lastDash := false
	for _, r := range lower {
		ok := r >= 'a' && r <= 'z' || r >= '0' && r <= '9'
		if ok {
			out = append(out, r)
			lastDash = false
			continue
		}
		if !lastDash {
			out = append(out, '-')
			lastDash = true
		}
	}
	slug := strings.Trim(string(out), "-")
	if slug == "" {
		return "child"
	}
	if len(slug) > 48 {
		slug = strings.Trim(slug[:48], "-")
	}
	return slug
}
