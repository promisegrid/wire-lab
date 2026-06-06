package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	poc13 "promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions"
)

func main() {
	configPath := flag.String("config", "config.json", "POC13 config path")
	containerName := flag.String("container", "", "container name from config")
	flag.Parse()
	if err := run(context.Background(), *configPath, *containerName); err != nil {
		fmt.Fprintf(os.Stderr, "poc13-supervisor: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, configPath, containerName string) error {
	if containerName == "" {
		return fmt.Errorf("container is required")
	}
	cfg, loadErr := poc13.LoadConfig(configPath)
	if loadErr != nil {
		return loadErr
	}
	container, ok := cfg.Container(containerName)
	if !ok {
		return fmt.Errorf("unknown container %s", containerName)
	}
	// Intent: The supervisor now owns the container-local TCP runtime so POC13
	// proves inter-container promise delivery instead of only per-agent log
	// generation. Source: DI-fumol
	runtime, runtimeErr := poc13.NewTCPRuntime(cfg, container, poc13.LiveOrScriptedDecider{})
	if runtimeErr != nil {
		return runtimeErr
	}
	return runtime.Run(ctx)
}
