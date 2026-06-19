package poc16_test

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestPayloadNamesDoNotUseObsoletePrefix(t *testing.T) {
	// Intent: POC16 payload names are pCID-local names. The old universal
	// prefix was bad vocabulary and must not return in code, tests, docs, or
	// retained examples. Source: DI-pusak
	forbiddenPrefix := "field" + "_"
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(entry.Name(), "core.") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if !utf8.Valid(data) {
			return nil
		}
		if bytes.Contains(data, []byte(forbiddenPrefix)) {
			t.Errorf("%s contains obsolete payload prefix", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
