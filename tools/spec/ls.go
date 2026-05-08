package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// cmdLs prints the manifest entries in a human-friendly table.
//
// Columns: pCID, slug, status, frozen-on, supersedes, freezing-commit.
// The depends-on and notes fields are omitted from the table for brevity;
// machine consumers should parse the YAML block in MANIFEST.md directly.
func cmdLs() error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}
	m, err := readManifestOrEmpty(repoRoot)
	if err != nil {
		return err
	}
	if len(m.Entries) == 0 {
		fmt.Println("(no frozen specs)")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(w, "PCID\tSLUG\tSTATUS\tFROZEN_ON\tSUPERSEDES\tCOMMIT")
	for _, e := range m.Entries {
		commit := e.FreezingCommit
		if len(commit) > 12 {
			commit = commit[:12]
		}
		sup := e.Supersedes
		if sup == "" {
			sup = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			e.PCID, e.Slug, e.Status, e.FrozenOn, sup, commit)
	}
	return w.Flush()
}
