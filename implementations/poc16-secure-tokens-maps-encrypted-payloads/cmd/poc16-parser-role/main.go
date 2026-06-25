package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/lifecycle"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/parserrole"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.json", "POC16 config path")
	containerName := flag.String("container", "", "container name from config")
	flag.Parse()
	if *containerName == "" {
		return fmt.Errorf("-container is required")
	}
	cfg, loadErr := config.Load(*configPath)
	if loadErr != nil {
		return loadErr
	}
	if _, containerFound := cfg.Container(*containerName); !containerFound {
		return fmt.Errorf("unknown container %q", *containerName)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// Intent: The parser role is an independent local kernel role, not a helper
	// function hidden inside app or transport-kernel code. Source: DI-gazin
	options := lifecycle.OptionsFromEnv("supervisor:"+*containerName, "parser:"+*containerName, lifecycle.RoleKindParser, cfg.RunID, pcid.NewRegistry().MustCID(pcid.LocalLifecycleV1), cfg.ShutdownGrace(), false, os.Stdin)
	if options.Address == "" {
		return parserrole.New(cfg, *containerName).Run(ctx)
	}
	managedRole, lifecycleErr := lifecycle.NewManagedRole(options)
	if lifecycleErr != nil {
		return lifecycleErr
	}
	role := parserrole.New(cfg, *containerName)
	role.LifecycleHandler = managedRole.HandleInvocationFrame
	return managedRole.Run(ctx, role.Run)
}
