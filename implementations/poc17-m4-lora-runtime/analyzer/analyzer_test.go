package analyzer

import (
	"os"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/config"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/sim"
)

func TestAnalyzeCleanRun(t *testing.T) {
	root := t.TempDir()
	cfg := config.Default()
	cfg.RunRoot = root
	cfg.RunID = "test-run"
	cfg.RadioMTUBytes = 96
	if err := sim.Run(cfg); err != nil {
		t.Fatal(err)
	}
	summary, err := Analyze(cfg.RunDir())
	if err != nil {
		t.Fatal(err)
	}
	if summary.RadioSends == 0 || summary.MessageArtifacts == 0 || summary.PeerStorage == 0 {
		t.Fatalf("incomplete summary: %+v", summary)
	}
	if _, err := os.Stat(cfg.RunDir()); err != nil {
		t.Fatal(err)
	}
}
