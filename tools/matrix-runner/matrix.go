package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func runUpdateMatrix(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("update-matrix", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	dryRun := fs.Bool("dry-run", false, "print intended updates without writing")
	var results multiFlag
	fs.Var(&results, "result", "result file path; repeatable")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(results) == 0 {
		return errUsage("update-matrix: --result is required")
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	for _, result := range results {
		matrix, action, row, err := updateMatrixForResult(repo, result, *dryRun)
		if err != nil {
			return err
		}
		if *dryRun {
			fmt.Fprintf(stdout, "%s: %s -> %s\n", action, repo.Rel(matrix), row)
		} else {
			fmt.Fprintf(stdout, "%s: %s\n", action, repo.Rel(matrix))
		}
	}
	return nil
}

type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func resultCoordinates(repo Repo, result string) (simID, scenarioID, modelID, timestamp, abs string, err error) {
	abs = repo.Abs(result)
	rel := repo.Rel(abs)
	parts := strings.Split(rel, "/")
	if len(parts) != 5 || parts[0] != "results" || !strings.HasSuffix(parts[4], ".md") {
		err = fmt.Errorf("result path must have shape results/<sim-id>/<scenario-id>/<model-id>/<timestamp>.md")
		return
	}
	return parts[1], parts[2], parts[3], strings.TrimSuffix(parts[4], ".md"), abs, nil
}

func updateMatrixForResult(repo Repo, result string, dryRun bool) (matrixPath string, action string, row string, err error) {
	simID, scenarioID, _, _, abs, err := resultCoordinates(repo, result)
	if err != nil {
		return "", "", "", err
	}
	verdict, err := extractVerdict(abs)
	if err != nil {
		return "", "", "", err
	}
	relResult := repo.Rel(abs)
	row = fmt.Sprintf("| `%s` | `%s` | `%s` | run | %s |", simID, scenarioID, relResult, markdownEscapeCell(verdict))
	matrixPath = repo.Path("scenarios", scenarioID, "MATRIX.md")
	bytes, err := os.ReadFile(matrixPath)
	if err != nil {
		return "", "", "", err
	}
	lines := strings.Split(strings.TrimRight(string(bytes), "\n"), "\n")
	replaced := false
	for index, line := range lines {
		cells := splitTableRow(line)
		if len(cells) < 2 {
			continue
		}
		if stripCodeCell(cells[0]) == simID && stripCodeCell(cells[1]) == scenarioID {
			lines[index] = row
			replaced = true
			break
		}
	}
	if replaced {
		action = "replace"
	} else {
		action = "append"
		lines = append(lines, row)
	}
	if !dryRun {
		if err := writeFile(matrixPath, joinLines(lines)); err != nil {
			return "", "", "", err
		}
	}
	return matrixPath, action, row, nil
}

func splitTableRow(line string) []string {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "|") {
		return nil
	}
	raw := strings.Split(strings.Trim(line, "|"), "|")
	cells := make([]string, 0, len(raw))
	for _, cell := range raw {
		cells = append(cells, strings.TrimSpace(cell))
	}
	return cells
}

func resultPathForCell(repo Repo, cell MatrixCell) string {
	return filepath.ToSlash(repo.Rel(repo.Abs(cell.ResultPath)))
}
