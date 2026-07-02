package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitAndDiscoverConfig(t *testing.T) {
	root := t.TempDir()
	repository, initErr := Init(root, "")
	if initErr != nil {
		t.Fatalf("Init() error = %v", initErr)
	}
	if repository.Config.CAS.Type != FileCASType || repository.Config.CAS.Path != DefaultCASPath {
		t.Fatalf("config CAS = %#v", repository.Config.CAS)
	}
	if _, statErr := os.Stat(filepath.Join(root, GridDirName, ConfigFileName)); statErr != nil {
		t.Fatalf("config missing: %v", statErr)
	}
	if _, statErr := os.Stat(filepath.Join(root, GridDirName, "cas", "objects")); statErr != nil {
		t.Fatalf("local CAS objects dir missing: %v", statErr)
	}
	nested := filepath.Join(root, "docs", "api")
	if mkdirErr := os.MkdirAll(nested, 0o755); mkdirErr != nil {
		t.Fatalf("MkdirAll() error = %v", mkdirErr)
	}
	discovered, discoverErr := Discover(nested)
	if discoverErr != nil {
		t.Fatalf("Discover() error = %v", discoverErr)
	}
	if discovered.Root != repository.Root {
		t.Fatalf("discovered root = %s, want %s", discovered.Root, repository.Root)
	}
}

func TestDiscoverReportsMissingConfig(t *testing.T) {
	root := t.TempDir()
	if _, discoverErr := Discover(root); discoverErr == nil {
		t.Fatalf("Discover() succeeded without .grid/config.json")
	}
}

func TestParseConfigRejectsUnimplementedCASType(t *testing.T) {
	_, parseErr := ParseConfig([]byte(`{"version":1,"cas":{"type":"daemon","path":"local"}}`))
	if parseErr == nil {
		t.Fatalf("ParseConfig() accepted unimplemented CAS type")
	}
}
