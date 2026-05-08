// Package main is the entry point for the `spec` command.
//
// `spec` is a single Go binary that handles the operational machinery of the
// PromiseGrid Wire Lab spec-doc store: freezing draft spec docs into
// content-addressed snapshots, auditing the manifest and on-disk state, and
// providing utility subcommands for humans and scripts.
//
// The binary's design is locked by TODO 011 (DI-011-...457 through ...501):
//
//   - Hash input is raw file bytes, no normalization (DI-011-20260429-184457).
//   - Tooling language is Go using github.com/ipfs/go-cid (DI-011-20260429-184458).
//   - Freezing produces a snapshot file and a manifest entry, no git tag
//     (DI-011-20260429-184459).
//   - Manifest is a single Markdown file with one fenced YAML block inside
//     (DI-011-20260429-184500).
//   - Freezer and checker are the SAME binary with subcommands (DI-011-20260429-184501).
//
// CIDv1 parameter set: multibase=base32, multihash=sha2-256, codec=raw.
package main

import (
	"fmt"
	"os"
)

const usage = `Usage:
  spec freeze <slug>      Mint a pCID, snapshot the draft, append manifest entry.
  spec check              Audit manifest, on-disk files, and cross-references.
  spec cid <file>         Print the CIDv1 of any file (utility).
  spec ls                 List frozen specs from the manifest.
  spec help               Print this message.

All paths are resolved relative to the repo root, which is auto-detected by
walking up from the current working directory until a '.git' directory is
found. The wire-lab harness's specs live under
'protocols/wire-lab.d/specs/' (SpecsDir) per TE-29 + TE-32.

Locked CIDv1 parameters (per DI-011-20260429-184453):
  multibase=base32, multihash=sha2-256, codec=raw
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	sub := os.Args[1]
	args := os.Args[2:]

	var err error
	switch sub {
	case "freeze":
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "spec freeze: exactly one <slug> argument required")
			os.Exit(2)
		}
		err = cmdFreeze(args[0])
	case "check":
		err = cmdCheck()
	case "cid":
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "spec cid: exactly one <file> argument required")
			os.Exit(2)
		}
		err = cmdCid(args[0])
	case "ls":
		err = cmdLs()
	case "help", "-h", "--help":
		fmt.Print(usage)
		return
	default:
		fmt.Fprintf(os.Stderr, "spec: unknown subcommand %q\n\n", sub)
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "spec %s: %v\n", sub, err)
		os.Exit(1)
	}
}
