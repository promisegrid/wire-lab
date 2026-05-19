package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PopulationSim describes one committed/tracked simulation available to GA
// selection.
//
// Intent: Keep ordinary GA population scans rooted in git-tracked simulation
// specimens, while generated untracked children stay pending until accepted.
// Source: DI-bagih
type PopulationSim struct {
	SimID    string
	Path     string
	Files    []string
	TreeHash string
}

func runInit(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	dryRun := fs.Bool("dry-run", false, "print tracked simulation population without writing GA state")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !*dryRun {
		return notImplemented("init")
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	population, err := discoverTrackedPopulation(repo)
	if err != nil {
		return err
	}
	if err := writeFormat(stdout, "population=%d\n", len(population)); err != nil {
		return err
	}
	for _, sim := range population {
		if err := writeFormat(stdout, "%s files=%d tree_hash=%s path=%s\n", sim.SimID, len(sim.Files), sim.TreeHash, sim.Path); err != nil {
			return err
		}
	}
	return nil
}

func discoverTrackedPopulation(repo Repo) ([]PopulationSim, error) {
	output, err := repo.Git("ls-files", "-z", "--", "simulations")
	if err != nil {
		return nil, err
	}
	bySim := map[string][]string{}
	for _, rel := range splitNUL(output) {
		simID, ok := simIDFromTrackedPath(rel)
		if !ok {
			continue
		}
		abs := repo.Abs(rel)
		info, err := os.Stat(abs)
		if err != nil || info.IsDir() {
			continue
		}
		bySim[simID] = append(bySim[simID], filepath.ToSlash(rel))
	}
	var simIDs []string
	for simID := range bySim {
		simIDs = append(simIDs, simID)
	}
	sort.Strings(simIDs)
	var population []PopulationSim
	for _, simID := range simIDs {
		files := bySim[simID]
		sort.Strings(files)
		treeHash, err := trackedTreeHash(repo, files)
		if err != nil {
			return nil, err
		}
		population = append(population, PopulationSim{
			SimID:    simID,
			Path:     filepath.ToSlash(filepath.Join("simulations", simID)) + "/",
			Files:    files,
			TreeHash: treeHash,
		})
	}
	return population, nil
}

func splitNUL(output []byte) []string {
	raw := bytes.Split(output, []byte{0})
	var values []string
	for _, item := range raw {
		value := strings.TrimSpace(string(item))
		if value != "" {
			values = append(values, filepath.ToSlash(value))
		}
	}
	return values
}

func simIDFromTrackedPath(rel string) (string, bool) {
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) < 3 || parts[0] != "simulations" || !strings.HasPrefix(parts[1], "SIM-") {
		return "", false
	}
	return parts[1], true
}

func trackedTreeHash(repo Repo, files []string) (string, error) {
	digest := sha256.New()
	for _, rel := range files {
		if err := hashWrite(digest, []byte(rel)); err != nil {
			return "", err
		}
		if err := hashWrite(digest, []byte{0}); err != nil {
			return "", err
		}
		content, err := os.ReadFile(repo.Abs(rel))
		if err != nil {
			return "", err
		}
		fileHash := sha256.Sum256(content)
		if err := hashWrite(digest, []byte(hex.EncodeToString(fileHash[:]))); err != nil {
			return "", err
		}
		if err := hashWrite(digest, []byte{0}); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func hashWrite(writer hash.Hash, data []byte) error {
	written, err := writer.Write(data)
	if err != nil {
		return err
	}
	if written != len(data) {
		return fmt.Errorf("hash write short count: wrote %d of %d", written, len(data))
	}
	return nil
}
