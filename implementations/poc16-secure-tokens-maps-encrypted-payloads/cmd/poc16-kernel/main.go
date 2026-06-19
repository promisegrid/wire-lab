package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/kernel"
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
	if _, ok := cfg.Container(*containerName); !ok {
		return fmt.Errorf("unknown container %q", *containerName)
	}
	// Intent: The kernel is a container-local transport process. It routes
	// exact envelopes between app processes and peer kernels without owning app
	// trust, workflow, or promise interpretation. Source: DI-galin
	localKernel := kernel.New(cfg, *containerName)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return localKernel.Run(ctx)
}
