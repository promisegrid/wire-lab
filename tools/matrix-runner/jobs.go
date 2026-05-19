package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
)

func runJobs(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("jobs", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	manifest := fs.String("manifest", "", "matrix manifest CSV path")
	outputDir := fs.String("output-dir", "", "prompt output directory")
	timestamp := fs.String("timestamp", timestampPlaceholder, "optional fixed timestamp override")
	maxCells := fs.Int("max-cells", -1, "optional maximum generated prompt count")
	startIndex := fs.Int("start-index", 0, "zero-based start index")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *manifest == "" {
		return errUsage("jobs: --manifest is required")
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	cells, err := readManifest(repo, *manifest)
	if err != nil {
		return err
	}
	for index := range cells {
		if *timestamp != "" && *timestamp != timestampPlaceholder {
			cells[index].Timestamp = *timestamp
			cells[index].ResultPath = defaultResultPath(cells[index].SimID, cells[index].ScenarioID, cells[index].ModelID, *timestamp)
		}
	}
	selected, err := selectedCells(cells, *startIndex, *maxCells)
	if err != nil {
		return err
	}
	outDir := *outputDir
	if outDir == "" {
		outDir = repo.Path("results", "jobs", selected[0].RunGroupID)
	} else {
		outDir = repo.Abs(outDir)
	}
	if err := osMkdirAll(outDir, 0o755); err != nil {
		return err
	}
	indexLines := []string{fmt.Sprintf("# LLM Jobs: %s", selected[0].RunGroupID), ""}
	for _, cell := range selected {
		name := promptFilename(cell)
		path := filepath.Join(outDir, name)
		if err := writeFile(path, pathOnlyPrompt(cell)); err != nil {
			return err
		}
		indexLines = append(indexLines, fmt.Sprintf("- [%s](%s)", name, name))
	}
	if err := writeFile(filepath.Join(outDir, "INDEX.md"), joinLines(indexLines)); err != nil {
		return err
	}
	if err := writeLine(stdout, repo.Rel(outDir)); err != nil {
		return err
	}
	return writeFormat(stdout, "jobs=%d\n", len(selected))
}
