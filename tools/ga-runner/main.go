// Package main implements ga-runner, the PromiseGrid Wire Lab tool for
// JSON-fitness GA/search runs.
//
// Intent: Start the GA/search runner on a clean contract instead of extending
// the Markdown-oriented matrix-runner path. Source: DI-pobus
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const usage = `Usage:
  ga-runner init       Create or preview a GA run state file.
  ga-runner backfill-init  Create a targeted rubric-v2 backfill state file.
  ga-runner compare-backfill  Compare targeted rubric-v2 backfill results against canonical v1 evidence.
  ga-runner score      Score GA cells and write JSON fitness results.
  ga-runner generate   Generate untracked child simulations.
  ga-runner validate   Validate JSON fitness result files.
  ga-runner audit      Audit canonical scored results for rubric-v2 backfill.
  ga-runner progress   Show GA run progress. (not implemented yet)
  ga-runner accept     Record accepted children and staging paths.
  ga-runner cull       Delete rejected generated children and their results.
  ga-runner help       Print this message.

Commands accept -repo-root. When omitted, ga-runner walks up from the current
directory until it finds a .git directory.
`

func main() {
	if err := runMain(os.Args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "ga-runner: %v\n", err)
		os.Exit(1)
	}
}

func runMain(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) < 2 {
		if err := writeText(stderr, usage); err != nil {
			return err
		}
		return exitCodeError{code: 2, message: "missing subcommand"}
	}
	subcommand := args[1]
	subArgs := args[2:]
	switch subcommand {
	case "validate":
		return runValidate(subArgs, stdout)
	case "init":
		return runInit(subArgs, stdout)
	case "backfill-init":
		// Intent: Build an audit-first rescore state file from historical v1
		// evidence without rewriting any scored artifact bytes. Source: DI-roruj
		return runBackfillInit(subArgs, stdout)
	case "compare-backfill":
		// Intent: Materialize a durable drift report from targeted rubric-v2
		// backfill evidence before operators broaden rescoring scope. Source:
		// DI-zuzup
		return runCompareBackfill(subArgs, stdout)
	case "accept":
		// Intent: Route review/promotion through the acceptance checkpoint
		// instead of the not-implemented stub. Source: DI-podot
		return runAccept(subArgs, stdout)
	case "cull":
		// Intent: Route destructive rejection cleanup through the state-bound cull
		// checkpoint instead of the not-implemented stub. Source: DI-kofil
		return runCull(subArgs, stdout)
	case "score":
		// Intent: Route JSON-fitness scoring through the stateful GA loop instead
		// of the legacy Markdown matrix contract. Source: DI-gijom
		return runScore(subArgs, stdout)
	case "generate":
		// Intent: Materialize child simulations only through state-bound GA
		// generation so pending children remain auditable. Source: DI-gijom
		return runGenerate(subArgs, stdout)
	case "audit":
		// Intent: Audit canonical GA evidence before a targeted vocabulary-aware
		// backfill so rescoring focuses on the sims most likely to move. Source:
		// DI-roruj
		return runAudit(subArgs, stdout)
	case "progress":
		return notImplemented(subcommand)
	case "help", "-h", "--help":
		return writeText(stdout, usage)
	default:
		if err := writeText(stderr, usage); err != nil {
			return err
		}
		return exitCodeError{code: 2, message: "unknown subcommand " + subcommand}
	}
}

func notImplemented(command string) error {
	return fmt.Errorf("%s: not implemented yet; see TODO-tapur", command)
}

type exitCodeError struct {
	code    int
	message string
}

func (e exitCodeError) Error() string {
	return e.message
}

func commonRepoFlag(fs *flag.FlagSet) *string {
	return fs.String("repo-root", "", "wire-lab repository root; default auto-detects by walking up to .git")
}

func errUsage(message string) error {
	return exitCodeError{code: 2, message: message}
}

func isExitCodeError(err error) bool {
	var target exitCodeError
	return errors.As(err, &target)
}
