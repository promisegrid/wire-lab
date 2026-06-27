package analyzer

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/config"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/protocol"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/sim"
)

func TestAnalyzeCleanRun(t *testing.T) {
	root := t.TempDir()
	cfg := config.Default()
	cfg.RunRoot = root
	cfg.RunID = "test-run"
	if err := sim.Run(cfg); err != nil {
		t.Fatal(err)
	}
	summary, err := Analyze(cfg.RunDir())
	if err != nil {
		t.Fatal(err)
	}
	if summary.RadioSends == 0 || summary.MessageArtifacts == 0 || summary.PeerStorageGrants == 0 || summary.PeerStorageGets == 0 || summary.DeviceRecoveries == 0 || summary.OrderStatusEvents == 0 || summary.LifecycleIssued == 0 || summary.ResourceWithdrawals == 0 {
		t.Fatalf("incomplete summary: %+v", summary)
	}
	if _, err := os.Stat(cfg.RunDir()); err != nil {
		t.Fatal(err)
	}
}

func TestCleanRunUsesCIDFirstOutput(t *testing.T) {
	root := t.TempDir()
	cfg := config.Default()
	cfg.RunRoot = root
	cfg.RunID = "cid-run"
	if err := sim.Run(cfg); err != nil {
		t.Fatal(err)
	}
	hexDigest := regexp.MustCompile(`^[a-f0-9]{64}$`)
	for _, dir := range []string{"message-cas", "malformed", "lifecycle-cas"} {
		entries, err := os.ReadDir(filepath.Join(cfg.RunDir(), dir))
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			stem := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			candidate := stem
			if strings.Contains(stem, "-") {
				parts := strings.Split(stem, "-")
				candidate = parts[len(parts)-1]
			}
			if hexDigest.MatchString(candidate) {
				t.Fatalf("%s uses bare SHA-256 digest %s", entry.Name(), candidate)
			}
			if err := protocol.ValidateCIDText(candidate); err != nil {
				t.Fatalf("%s is not canonical CID text: %v", entry.Name(), err)
			}
		}
	}
	f, err := os.Open(filepath.Join(cfg.RunDir(), "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			t.Fatalf("close events: %v", closeErr)
		}
	}()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), `"hash"`) || strings.Contains(scanner.Text(), "missing-parent-cid") {
			t.Fatalf("non-CID-first event line: %s", scanner.Text())
		}
		var event artifact.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			t.Fatal(err)
		}
		for _, value := range []string{event.PCID, event.CID} {
			if value == "" {
				continue
			}
			if err := protocol.ValidateCIDText(value); err != nil {
				t.Fatalf("event %s has non-canonical CID %q: %v", event.Type, value, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}
