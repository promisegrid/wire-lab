package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"strconv"
)

func runManifest(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("manifest", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	models := fs.String("models", "", "comma-separated model IDs")
	runGroupID := fs.String("run-group-id", "", "run group ID; default current UTC timestamp")
	timestamp := fs.String("timestamp", "", "result timestamp; default run group timestamp or current UTC")
	output := fs.String("output", "", "output CSV path")
	simGlob := fs.String("sim-glob", "SIM-*", "simulation directory glob under simulations/")
	shuffleSeed := fs.String("shuffle-seed", "", "optional deterministic shuffle seed")
	limitCells := fs.Int("limit-cells", -1, "optional maximum emitted rows")
	if err := fs.Parse(args); err != nil {
		return err
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	modelIDs := splitCommaList(*models)
	if len(modelIDs) == 0 {
		return errUsage("manifest: --models is required")
	}
	generated := utcCompactTimestamp()
	group := *runGroupID
	if group == "" {
		group = generated
	}
	ts := *timestamp
	if ts == "" {
		if *runGroupID == "" {
			ts = group
		} else {
			ts = generated
		}
	}
	simIDs, err := discoverSimIDs(repo, *simGlob)
	if err != nil {
		return err
	}
	if len(simIDs) == 0 {
		return fmt.Errorf("no simulations matched glob %q", *simGlob)
	}
	scenarioIDs, err := discoverScenarioIDs(repo)
	if err != nil {
		return err
	}
	if len(scenarioIDs) == 0 {
		return fmt.Errorf("no scenario entries found under scenarios/")
	}
	cells := buildManifestRows(group, ts, simIDs, scenarioIDs, modelIDs)
	if *shuffleSeed != "" {
		seed, err := strconv.ParseInt(*shuffleSeed, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid --shuffle-seed: %w", err)
		}
		rand.New(rand.NewSource(seed)).Shuffle(len(cells), func(i, j int) {
			cells[i], cells[j] = cells[j], cells[i]
		})
		assignOrdinals(cells)
	}
	if *limitCells >= 0 && *limitCells < len(cells) {
		cells = cells[:*limitCells]
		assignOrdinals(cells)
	}
	outPath := *output
	if outPath == "" {
		outPath = filepath.Join(repo.Root, "results", "manifests", "matrix-manifest-"+group+".csv")
	} else {
		outPath = repo.Abs(outPath)
	}
	if err := writeManifest(outPath, cells); err != nil {
		return err
	}
	fmt.Fprintln(stdout, repo.Rel(outPath))
	fmt.Fprintf(stdout, "rows=%d sims=%d scenarios=%d models=%d timestamp=%s\n", len(cells), len(simIDs), len(scenarioIDs), len(modelIDs), ts)
	return nil
}
