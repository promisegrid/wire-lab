package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func pathOnlyPrompt(cell MatrixCell, resultStyle string) string {
	var out strings.Builder
	fmt.Fprintln(&out, "# LLM Matrix Cell Job")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "## Cell")
	fmt.Fprintln(&out)
	fmt.Fprintf(&out, "- Cell ID: `%s`\n", cell.CellID)
	fmt.Fprintf(&out, "- Run group ID: `%s`\n", cell.RunGroupID)
	fmt.Fprintf(&out, "- Queue ordinal: `%d`\n", cell.Ordinal)
	fmt.Fprintf(&out, "- Simulation ID: `%s`\n", cell.SimID)
	fmt.Fprintf(&out, "- Scenario ID: `%s`\n", cell.ScenarioID)
	fmt.Fprintf(&out, "- Model ID: `%s`\n", cell.ModelID)
	fmt.Fprintf(&out, "- Intended result path: `%s`\n", cell.ResultPath)
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "## Required Source Inputs")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "Read only source/design inputs before producing the verdict:")
	fmt.Fprintln(&out)
	fmt.Fprintf(&out, "- `%sREADME.md`\n", cell.SimPath)
	fmt.Fprintf(&out, "- `%sQUESTION.md` if present\n", cell.SimPath)
	fmt.Fprintf(&out, "- local draft specs under `%s` if present\n", cell.SimPath)
	fmt.Fprintln(&out, "- `scenarios/README.md`")
	fmt.Fprintf(&out, "- `%s`\n", cell.ScenarioPath)
	fmt.Fprintf(&out, "- local scenario docs under `scenarios/%s/` if present\n", cell.ScenarioID)
	fmt.Fprintln(&out, "- `results/RUN-PROTOCOL.md`")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "Do not read prior result files for this same sim/scenario cell before writing")
	fmt.Fprintln(&out, "the verdict. This job is blind with respect to prior results.")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "## Task")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "Evaluate the simulation against the scenario using deeper reasoning. Explain:")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "- what the simulation can actually cover,")
	fmt.Fprintln(&out, "- what obligations it pushes to another layer,")
	fmt.Fprintln(&out, "- where the scenario's 100-year, sparse-knowledge, no-central-authority,")
	fmt.Fprintln(&out, "  auditability, and migration pressures expose weaknesses,")
	fmt.Fprintln(&out, "- which open questions remain.")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "## Result Style")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, resultStyleInstruction(resultStyle))
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "Write the result file at:")
	fmt.Fprintln(&out)
	fmt.Fprintf(&out, "`%s`\n", cell.ResultPath)
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "The result must follow the section contract in `results/RUN-PROTOCOL.md` and")
	fmt.Fprintln(&out, "must include:")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "- `- Run mode: llm-doc-eval-blind`")
	fmt.Fprintln(&out, "- a line starting with `Evidence verdict:`")
	fmt.Fprintln(&out, "- an explicit `Authority Boundary` section.")
	return out.String()
}

type SourceDocument struct {
	Path string
	Text string
}

type PromptBuilder struct {
	Repo        Repo
	ResultStyle string
}

func (b PromptBuilder) BuildAPIPrompt(cell MatrixCell) (string, error) {
	docs, err := b.sourceDocuments(cell)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	fmt.Fprintf(&out, "# Matrix Cell API Evaluation\n\n")
	fmt.Fprintf(&out, "Return only the complete Markdown result file. Do not wrap it in a code fence. Do not ask for confirmation.\n\n")
	fmt.Fprintf(&out, "## Result Style\n\n")
	fmt.Fprintf(&out, "%s\n\n", resultStyleInstruction(b.ResultStyle))
	fmt.Fprintf(&out, "## Required Result Contract\n\n")
	fmt.Fprintf(&out, "The Markdown result must contain exactly the required sections from `results/RUN-PROTOCOL.md`, include `- Run mode: llm-doc-eval-blind`, include one line starting with `Evidence verdict:`, and include an explicit `Authority Boundary` section.\n\n")
	fmt.Fprintf(&out, "## Evaluation Task\n\n")
	fmt.Fprintf(&out, "Evaluate the simulation against the scenario. Explain what the simulation covers, what it pushes to another layer, how the scenario's 100-year durability, sparse knowledge, no central authority, auditability, and migration pressures expose weaknesses, and which open questions remain.\n\n")
	fmt.Fprintf(&out, "## Required Output Coordinates\n\n")
	fmt.Fprintf(&out, "- Result path: `%s`\n", cell.ResultPath)
	fmt.Fprintf(&out, "- Header: `# Result: %s / %s / %s / %s`\n", cell.SimID, cell.ScenarioID, cell.ModelID, cell.Timestamp)
	fmt.Fprintf(&out, "- Result ID: `%s-%s-%s-%s`\n", cell.SimID, cell.ScenarioID, cell.ModelID, cell.Timestamp)
	fmt.Fprintf(&out, "- Scenario ID: `%s`\n", cell.ScenarioID)
	fmt.Fprintf(&out, "- Scenario path: `%s`\n", cell.ScenarioPath)
	fmt.Fprintf(&out, "- Simulation ID: `%s`\n", cell.SimID)
	fmt.Fprintf(&out, "- Simulation path: `%s`\n", cell.SimPath)
	fmt.Fprintf(&out, "- Simulation commit: `%s`\n", b.Repo.GitCommit())
	fmt.Fprintf(&out, "- Model ID: `%s`\n", cell.ModelID)
	fmt.Fprintf(&out, "- Run timestamp UTC: `%s`\n\n", cell.Timestamp)
	fmt.Fprintf(&out, "## Source Documents\n\n")
	for _, doc := range docs {
		fmt.Fprintf(&out, "### `%s`\n\n", doc.Path)
		fmt.Fprintf(&out, "```markdown\n%s\n```\n\n", strings.TrimSpace(doc.Text))
	}
	return out.String(), nil
}

func validResultStyle(style string) bool {
	return style == "" || style == "concise" || style == "standard"
}

func resultStyleInstruction(style string) string {
	if style == "standard" {
		return "Use standard evidence prose. Keep the result focused, but include enough detail to audit the fit."
	}
	// Intent: Default API jobs to shorter, audit-focused results so full-matrix
	// runs spend fewer output tokens while preserving comparison evidence.
	// Source: DI-nugiv
	return "Use concise evidence prose. Keep each narrative section to 1-3 bullets or short paragraphs, avoid restating source documents, target roughly 700-1100 visible words, and prefer specific fit/weakness/open-question bullets over long exposition. Source: DI-nugiv"
}

func (b PromptBuilder) sourceDocuments(cell MatrixCell) ([]SourceDocument, error) {
	required := []string{
		"results/RUN-PROTOCOL.md",
		"scenarios/README.md",
		filepath.ToSlash(filepath.Join("scenarios", cell.ScenarioID, "README.md")),
		cell.ScenarioPath,
		filepath.ToSlash(filepath.Join("simulations", cell.SimID, "README.md")),
	}
	var docs []SourceDocument
	for _, rel := range required {
		text, err := b.Repo.ReadRel(rel)
		if err != nil {
			return nil, fmt.Errorf("read required source %s: %w", rel, err)
		}
		docs = append(docs, SourceDocument{Path: rel, Text: text})
	}
	questionPath := filepath.ToSlash(filepath.Join("simulations", cell.SimID, "QUESTION.md"))
	if text, err := b.Repo.ReadRel(questionPath); err == nil {
		docs = append(docs, SourceDocument{Path: questionPath, Text: text})
	}
	scenarioDocs, err := b.localScenarioMarkdown(cell)
	if err != nil {
		return nil, err
	}
	docs = append(docs, scenarioDocs...)
	localDocs, err := b.localSimulationMarkdown(cell)
	if err != nil {
		return nil, err
	}
	docs = append(docs, localDocs...)
	return docs, nil
}

// localScenarioMarkdown collects scenario-local markdown files that are not
// already included as required sources. It explicitly ignores legacy MATRIX.md
// names so stale generated or untracked summaries cannot leak into blind cell
// prompts. Source: DI-zamin
func (b PromptBuilder) localScenarioMarkdown(cell MatrixCell) ([]SourceDocument, error) {
	root := b.Repo.Path("scenarios", cell.ScenarioID)
	mainName := filepath.Base(cell.ScenarioPath)
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		base := filepath.Base(path)
		if base == "README.md" || base == mainName || base == "MATRIX.md" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return b.readSourceDocuments(paths)
}

// localSimulationMarkdown collects simulation-local design documents beyond the
// required README and QUESTION files.
func (b PromptBuilder) localSimulationMarkdown(cell MatrixCell) ([]SourceDocument, error) {
	root := b.Repo.Path("simulations", cell.SimID)
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		base := filepath.Base(path)
		if base == "README.md" || base == "QUESTION.md" || base == "SCENARIOS.md" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return b.readSourceDocuments(paths)
}

// readSourceDocuments reads source files in stable order so generated prompts
// are deterministic across equivalent filesystem walks.
func (b PromptBuilder) readSourceDocuments(paths []string) ([]SourceDocument, error) {
	sort.Strings(paths)
	var docs []SourceDocument
	for _, path := range paths {
		rel := b.Repo.Rel(path)
		text, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		docs = append(docs, SourceDocument{Path: rel, Text: string(text)})
	}
	return docs, nil
}
