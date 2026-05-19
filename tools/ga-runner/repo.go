package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Repo keeps all path decisions rooted in one discovered wire-lab checkout.
//
// Intent: Keep validation from accidentally accepting paths outside the repo or
// depending on the operator's current working directory. Source: DI-pobus
type Repo struct {
	Root string
}

func openRepo(root string) (Repo, error) {
	if root != "" {
		abs, err := filepath.Abs(root)
		if err != nil {
			return Repo{}, err
		}
		return Repo{Root: abs}, nil
	}
	working, err := os.Getwd()
	if err != nil {
		return Repo{}, err
	}
	for {
		if info, err := os.Stat(filepath.Join(working, ".git")); err == nil && info.IsDir() {
			return Repo{Root: working}, nil
		}
		parent := filepath.Dir(working)
		if parent == working {
			return Repo{}, fmt.Errorf("could not find repository root from %s", working)
		}
		working = parent
	}
}

func (repo Repo) Path(parts ...string) string {
	all := append([]string{repo.Root}, parts...)
	return filepath.Join(all...)
}

func (repo Repo) Abs(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(repo.Root, path)
}

func (repo Repo) Rel(path string) string {
	rel, err := filepath.Rel(repo.Root, path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}
	return filepath.ToSlash(rel)
}

func (repo Repo) Git(args ...string) ([]byte, error) {
	// Intent: Query git's index for stable population membership so ordinary
	// scans are not contaminated by untracked generated child sims. Source:
	// DI-bagih
	fullArgs := append([]string{"-C", repo.Root}, args...)
	command := exec.Command("git", fullArgs...)
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("git %v: %w", args, err)
	}
	return output, nil
}
