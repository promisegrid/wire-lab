package main

import (
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const timestampPlaceholder = "<YYYYMMDD-HHMMSS>"

var slugPattern = regexp.MustCompile(`[^A-Za-z0-9_.-]+`)

// MatrixCell is one concrete simulation/scenario/model work unit.
//
// Intent: Keep manifest, prompt, queue, validation, and matrix-update logic on
// one normalized data shape so unattended runs cannot drift between inferred
// result paths and committed result paths. Source: DI-lulom
type MatrixCell struct {
	RunGroupID     string
	Ordinal        int
	CellID         string
	SimID          string
	ScenarioID     string
	ModelID        string
	SimPath        string
	ScenarioPath   string
	ResultDir      string
	ResultPathTmpl string
	Timestamp      string
	ResultPath     string
	Status         string
}

func utcCompactTimestamp() string {
	return time.Now().UTC().Format("20060102-150405")
}

func slug(value string) string {
	value = slugPattern.ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}

func defaultResultDir(simID, scenarioID, modelID string) string {
	return fmt.Sprintf("results/%s/%s/%s/", simID, scenarioID, modelID)
}

func defaultResultPath(simID, scenarioID, modelID, timestamp string) string {
	return defaultResultDir(simID, scenarioID, modelID) + timestamp + ".md"
}

func defaultResultTemplate(simID, scenarioID, modelID string) string {
	return defaultResultPath(simID, scenarioID, modelID, timestampPlaceholder)
}

func defaultCellID(runGroupID string, ordinal int, simID, scenarioID, modelID string) string {
	return fmt.Sprintf("%s-%06d-%s--%s--%s", slug(runGroupID), ordinal, slug(simID), slug(scenarioID), slug(modelID))
}

func promptFilename(cell MatrixCell) string {
	return fmt.Sprintf("%05d-%s--%s--%s.md", cell.Ordinal, slug(cell.SimID), slug(cell.ScenarioID), slug(cell.ModelID))
}

func discoverSimIDs(repo Repo, simGlob string) ([]string, error) {
	pattern := simGlob
	if !filepath.IsAbs(pattern) {
		pattern = repo.Path("simulations", simGlob)
	}
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, match := range matches {
		info, err := osStat(match)
		if err == nil && info.IsDir() && strings.HasPrefix(filepath.Base(match), "SIM-") {
			ids = append(ids, filepath.Base(match))
		}
	}
	sort.Strings(ids)
	return ids, nil
}

func discoverScenarioIDs(repo Repo) ([]string, error) {
	entries, err := osReadDir(repo.Path("scenarios"))
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		if _, err := osStat(repo.Path("scenarios", id, id+".md")); err == nil {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids, nil
}

func buildManifestRows(runGroupID, timestamp string, simIDs, scenarioIDs, modelIDs []string) []MatrixCell {
	var cells []MatrixCell
	for _, simID := range simIDs {
		for _, scenarioID := range scenarioIDs {
			for _, modelID := range modelIDs {
				cells = append(cells, MatrixCell{
					RunGroupID:     runGroupID,
					SimID:          simID,
					ScenarioID:     scenarioID,
					ModelID:        modelID,
					SimPath:        fmt.Sprintf("simulations/%s/", simID),
					ScenarioPath:   fmt.Sprintf("scenarios/%s/%s.md", scenarioID, scenarioID),
					ResultDir:      defaultResultDir(simID, scenarioID, modelID),
					ResultPathTmpl: defaultResultTemplate(simID, scenarioID, modelID),
					Timestamp:      timestamp,
					ResultPath:     defaultResultPath(simID, scenarioID, modelID, timestamp),
					Status:         "queued",
				})
			}
		}
	}
	assignOrdinals(cells)
	return cells
}

func assignOrdinals(cells []MatrixCell) {
	if len(cells) == 0 {
		return
	}
	runGroupID := cells[0].RunGroupID
	for i := range cells {
		cells[i].Ordinal = i + 1
		cells[i].CellID = defaultCellID(runGroupID, cells[i].Ordinal, cells[i].SimID, cells[i].ScenarioID, cells[i].ModelID)
	}
}

func readManifest(repo Repo, path string) ([]MatrixCell, error) {
	handle, err := osOpen(repo.Abs(path))
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	reader := csv.NewReader(handle)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("manifest has no header")
	}
	header := map[string]int{}
	for index, name := range rows[0] {
		header[name] = index
	}
	required := []string{"run_group_id", "sim_id", "scenario_id", "model_id", "sim_path", "scenario_path"}
	for _, name := range required {
		if _, ok := header[name]; !ok {
			return nil, fmt.Errorf("manifest missing required field %q", name)
		}
	}
	var cells []MatrixCell
	for rowIndex, row := range rows[1:] {
		get := func(name string) string {
			index, ok := header[name]
			if !ok || index >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[index])
		}
		ordinal, _ := strconv.Atoi(get("ordinal"))
		if ordinal == 0 {
			ordinal = rowIndex + 1
		}
		runGroupID := get("run_group_id")
		simID := get("sim_id")
		scenarioID := get("scenario_id")
		modelID := get("model_id")
		timestamp := get("timestamp")
		resultDir := get("result_dir")
		if resultDir == "" {
			resultDir = defaultResultDir(simID, scenarioID, modelID)
		}
		resultTemplate := get("result_path_template")
		if resultTemplate == "" {
			resultTemplate = defaultResultTemplate(simID, scenarioID, modelID)
		}
		resultPath := get("result_path")
		if resultPath == "" && timestamp != "" {
			resultPath = strings.ReplaceAll(resultTemplate, timestampPlaceholder, timestamp)
		}
		if resultPath == "" {
			resultPath = defaultResultPath(simID, scenarioID, modelID, timestampPlaceholder)
		}
		cellID := get("cell_id")
		if cellID == "" {
			cellID = defaultCellID(runGroupID, ordinal, simID, scenarioID, modelID)
		}
		status := get("status")
		if status == "" {
			status = "queued"
		}
		cells = append(cells, MatrixCell{
			RunGroupID:     runGroupID,
			Ordinal:        ordinal,
			CellID:         cellID,
			SimID:          simID,
			ScenarioID:     scenarioID,
			ModelID:        modelID,
			SimPath:        get("sim_path"),
			ScenarioPath:   get("scenario_path"),
			ResultDir:      resultDir,
			ResultPathTmpl: resultTemplate,
			Timestamp:      timestamp,
			ResultPath:     resultPath,
			Status:         status,
		})
	}
	if len(cells) == 0 {
		return nil, fmt.Errorf("manifest contains no rows")
	}
	return cells, nil
}

func writeManifest(path string, cells []MatrixCell) error {
	if err := ensureParent(path); err != nil {
		return err
	}
	handle, err := osCreate(path)
	if err != nil {
		return err
	}
	defer handle.Close()
	writer := csv.NewWriter(handle)
	header := []string{
		"run_group_id", "ordinal", "cell_id", "sim_id", "scenario_id", "model_id",
		"sim_path", "scenario_path", "result_dir", "result_path_template",
		"timestamp", "result_path", "status",
	}
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, cell := range cells {
		row := []string{
			cell.RunGroupID,
			strconv.Itoa(cell.Ordinal),
			cell.CellID,
			cell.SimID,
			cell.ScenarioID,
			cell.ModelID,
			cell.SimPath,
			cell.ScenarioPath,
			cell.ResultDir,
			cell.ResultPathTmpl,
			cell.Timestamp,
			cell.ResultPath,
			cell.Status,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func selectedCells(cells []MatrixCell, startIndex int, limit int) ([]MatrixCell, error) {
	if startIndex < 0 {
		return nil, fmt.Errorf("start index must be non-negative")
	}
	if startIndex > len(cells) {
		return nil, fmt.Errorf("start index %d beyond %d cells", startIndex, len(cells))
	}
	selected := cells[startIndex:]
	if limit >= 0 && limit < len(selected) {
		selected = selected[:limit]
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("no cells selected after filters")
	}
	return selected, nil
}

func contentHashID(parts ...string) string {
	sum := sha1.Sum([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:])[:12]
}
