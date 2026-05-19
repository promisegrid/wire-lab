// Package main implements matrix-runner, the durable Go CLI for PromiseGrid
// Wire Lab result-matrix runs.
//
// Intent: Replace the Python matrix orchestration scripts with one typed Go
// program that can generate manifests, run unattended API-backed cells, resume
// from checkpoints, validate result evidence, generate read-only result views,
// and compare model corpora. Source: DI-lulom; DI-zamin
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const usage = `Usage:
  matrix-runner manifest       Generate a scenario/simulation/model manifest CSV.
  matrix-runner jobs           Generate blind LLM job prompts from a manifest.
  matrix-runner run            Run a manifest through an API provider with checkpoints.
  matrix-runner progress       Show queue checkpoint counts.
  matrix-runner validate       Validate result files or manifest rows.
  matrix-runner view           Generate a result view from canonical result files.
  matrix-runner compare        Compare verdict drift between two model corpora.
  matrix-runner help           Print this message.

All commands accept -repo-root. When omitted, matrix-runner walks up from the
current directory until it finds a .git directory.
`

func main() {
	if err := runMain(context.Background(), os.Args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "matrix-runner: %v\n", err)
		os.Exit(1)
	}
}

func runMain(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) < 2 {
		if err := writeText(stderr, usage); err != nil {
			return err
		}
		return exitCodeError{code: 2, message: "missing subcommand"}
	}
	subcommand := args[1]
	subArgs := args[2:]
	switch subcommand {
	case "manifest":
		return runManifest(subArgs, stdout)
	case "jobs":
		return runJobs(subArgs, stdout)
	case "run":
		return runQueue(ctx, subArgs, stdout)
	case "progress":
		return runProgress(subArgs, stdout)
	case "validate":
		return runValidate(subArgs, stdout)
	case "view":
		return runView(subArgs, stdout)
	case "compare":
		return runCompare(subArgs, stdout)
	case "help", "-h", "--help":
		return writeText(stdout, usage)
	default:
		if err := writeText(stderr, usage); err != nil {
			return err
		}
		return exitCodeError{code: 2, message: "unknown subcommand " + subcommand}
	}
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

func splitCommaList(value string) []string {
	var out []string
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func errUsage(message string) error {
	return exitCodeError{code: 2, message: message}
}

func isExitCodeError(err error) bool {
	var target exitCodeError
	return errors.As(err, &target)
}
