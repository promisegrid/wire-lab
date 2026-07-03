// Package repo owns the local `.grid/config.json` control file for POC18.
//
// Intent: Keep repo-local command discovery separate from the CAS itself. The
// config may point at `.grid/cas` today, but later it can point at a local daemon
// or remote CAS without changing normal CLI discovery. Source: DI-pahor
package repo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

const (
	// GridDirName is the repo-local PromiseGrid control directory.
	GridDirName = ".grid"
	// ConfigFileName is the human-editable repo config file inside `.grid`.
	ConfigFileName = "config.json"
	// StateFileName is the mutable local CLI state file inside `.grid`.
	StateFileName = "state.json"
	// DefaultCASPath is the first POC18 file-CAS locator written by `grid init`.
	DefaultCASPath = ".grid/cas"
	// FileCASType is the only CAS locator type implemented in this POC18 slice.
	FileCASType = "file"
)

// Config is the `.grid/config.json` shape.
type Config struct {
	Version int        `json:"version"`
	CAS     CASLocator `json:"cas"`
}

// CASLocator names where this repo's sparse CAS is found.
type CASLocator struct {
	Type string `json:"type"`
	Path string `json:"path,omitempty"`
}

// Repository records a discovered or initialized repo control directory.
type Repository struct {
	Root       string
	GridDir    string
	ConfigPath string
	Config     Config
}

// Init creates `.grid/config.json` under root and opens the configured file CAS.
func Init(root string, casPath string) (Repository, error) {
	absRoot, rootErr := filepath.Abs(root)
	if rootErr != nil {
		return Repository{}, rootErr
	}
	if casPath == "" {
		casPath = DefaultCASPath
	}
	repository := Repository{
		Root:       absRoot,
		GridDir:    filepath.Join(absRoot, GridDirName),
		ConfigPath: filepath.Join(absRoot, GridDirName, ConfigFileName),
		Config: Config{
			Version: 1,
			CAS:     CASLocator{Type: FileCASType, Path: casPath},
		},
	}
	if _, statErr := os.Stat(repository.ConfigPath); statErr == nil {
		return Repository{}, fmt.Errorf("%s already exists", repository.ConfigPath)
	} else if !os.IsNotExist(statErr) {
		return Repository{}, statErr
	}
	if err := os.MkdirAll(repository.GridDir, 0o755); err != nil {
		return Repository{}, err
	}
	if _, openErr := repository.OpenFileCAS(); openErr != nil {
		return Repository{}, openErr
	}
	content, marshalErr := json.MarshalIndent(repository.Config, "", "  ")
	if marshalErr != nil {
		return Repository{}, marshalErr
	}
	content = append(content, '\n')
	if writeErr := os.WriteFile(repository.ConfigPath, content, 0o644); writeErr != nil {
		return Repository{}, writeErr
	}
	if stateErr := repository.SaveState(defaultState()); stateErr != nil {
		return Repository{}, stateErr
	}
	return repository, nil
}

// Discover walks upward from start until it finds `.grid/config.json`.
func Discover(start string) (Repository, error) {
	absStart, absErr := filepath.Abs(start)
	if absErr != nil {
		return Repository{}, absErr
	}
	info, statErr := os.Stat(absStart)
	if statErr != nil {
		return Repository{}, statErr
	}
	current := absStart
	if !info.IsDir() {
		current = filepath.Dir(absStart)
	}
	for {
		configPath := filepath.Join(current, GridDirName, ConfigFileName)
		content, readErr := os.ReadFile(configPath)
		if readErr == nil {
			config, parseErr := ParseConfig(content)
			if parseErr != nil {
				return Repository{}, fmt.Errorf("%s: %w", configPath, parseErr)
			}
			return Repository{
				Root:       current,
				GridDir:    filepath.Join(current, GridDirName),
				ConfigPath: configPath,
				Config:     config,
			}, nil
		}
		if !os.IsNotExist(readErr) {
			return Repository{}, readErr
		}
		parent := filepath.Dir(current)
		if parent == current {
			return Repository{}, fmt.Errorf("no %s found from %s upward", filepath.Join(GridDirName, ConfigFileName), absStart)
		}
		current = parent
	}
}

// ParseConfig validates the `.grid/config.json` shape.
func ParseConfig(content []byte) (Config, error) {
	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return Config{}, err
	}
	if config.Version != 1 {
		return Config{}, fmt.Errorf("unsupported config version %d", config.Version)
	}
	if config.CAS.Type == "" {
		return Config{}, fmt.Errorf("cas.type is required")
	}
	if config.CAS.Type != FileCASType {
		return Config{}, fmt.Errorf("cas.type %q is not implemented in POC18", config.CAS.Type)
	}
	if config.CAS.Path == "" {
		return Config{}, fmt.Errorf("cas.path is required for file CAS")
	}
	return config, nil
}

// OpenFileCAS opens the configured file CAS.
func (repository Repository) OpenFileCAS() (*store.FileStore, error) {
	if repository.Config.CAS.Type != FileCASType {
		return nil, fmt.Errorf("cas.type %q is not implemented in POC18", repository.Config.CAS.Type)
	}
	return store.Open(repository.ResolvePath(repository.Config.CAS.Path))
}

// ResolvePath resolves relative repo config paths against the repo root.
func (repository Repository) ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(repository.Root, path)
}
