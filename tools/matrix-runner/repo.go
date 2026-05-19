package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Repo centralizes repository-relative path handling.
//
// Intent: Matrix tooling should be runnable from any working directory without
// changing the committed result path shape or accidentally writing outside the
// wire-lab checkout. Source: DI-lulom
type Repo struct {
	Root string
}

func openRepo(rootFlag string) (Repo, error) {
	if rootFlag != "" {
		root, err := filepath.Abs(rootFlag)
		if err != nil {
			return Repo{}, err
		}
		return Repo{Root: root}, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return Repo{}, err
	}
	for {
		if info, err := os.Stat(filepath.Join(wd, ".git")); err == nil && info.IsDir() {
			return Repo{Root: wd}, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return Repo{}, fmt.Errorf("could not find repo root containing .git")
		}
		wd = parent
	}
}

func (r Repo) Path(parts ...string) string {
	all := append([]string{r.Root}, parts...)
	return filepath.Join(all...)
}

func (r Repo) Abs(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(r.Root, path)
}

func (r Repo) Rel(path string) string {
	abs := r.Abs(path)
	rel, err := filepath.Rel(r.Root, abs)
	if err != nil {
		return filepath.ToSlash(path)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(abs)
	}
	return filepath.ToSlash(rel)
}

func (r Repo) ReadRel(path string) (string, error) {
	bytes, err := os.ReadFile(r.Abs(path))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (r Repo) GitCommit() string {
	cmd := exec.Command("git", "-C", r.Root, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	commit := strings.TrimSpace(string(out))
	if commit == "" {
		return "unknown"
	}
	return commit
}

func ensureParent(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0o755)
}

func writeFile(path string, data string) error {
	if err := ensureParent(path); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(data), 0o644)
}
