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
  ga-runner init       Create a GA run state file. (not implemented yet)
  ga-runner score      Score GA cells and write JSON fitness results. (not implemented yet)
  ga-runner generate   Generate untracked child simulations. (not implemented yet)
  ga-runner validate   Validate JSON fitness result files.
  ga-runner progress   Show GA run progress. (not implemented yet)
  ga-runner accept     Record accepted children and staging paths.
  ga-runner cull       Delete rejected generated children and their results.
  ga-runner help       Print this message.

Implemented commands accept -repo-root. When omitted, ga-runner walks up from the
current directory until it finds a .git directory.
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
	case "accept":
		// Intent: Route review/promotion through the acceptance checkpoint
		// instead of the not-implemented stub. Source: DI-podot
		return runAccept(subArgs, stdout)
	case "cull":
		// Intent: Route destructive rejection cleanup through the state-bound cull
		// checkpoint instead of the not-implemented stub. Source: DI-kofil
		return runCull(subArgs, stdout)
	case "score", "generate", "progress":
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
