package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.json", "POC11 config path")
	containerName := flag.String("container", "", "container name from config")
	nodeBinary := flag.String("node-bin", "poc11-node", "poc11 node binary path")
	flag.Parse()
	if *containerName == "" {
		return fmt.Errorf("-container is required")
	}
	cfg, loadErr := config.Load(*configPath)
	if loadErr != nil {
		return loadErr
	}
	container, containerFound := cfg.Container(*containerName)
	if !containerFound {
		return fmt.Errorf("unknown container %q", *containerName)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return runContainerAgents(ctx, *nodeBinary, *configPath, container)
}

func runContainerAgents(ctx context.Context, nodeBinary, configPath string, container config.ContainerConfig) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errs := make(chan error, len(container.Agents))
	var waitGroup sync.WaitGroup
	for _, agentName := range container.Agents {
		waitGroup.Add(1)
		go func(localAgentName string) {
			defer waitGroup.Done()
			errs <- runAgentProcess(ctx, nodeBinary, configPath, localAgentName)
		}(agentName)
	}
	go func() {
		waitGroup.Wait()
		close(errs)
	}()
	var firstErr error
	for agentErr := range errs {
		if agentErr != nil && firstErr == nil {
			firstErr = agentErr
			cancel()
		}
	}
	return firstErr
}

func runAgentProcess(ctx context.Context, nodeBinary, configPath, agentName string) error {
	// Intent: A container supervisor starts multiple independent agent kernels
	// without sharing their decision state; stdout/stderr are passed through so
	// the run log remains auditable from Docker output. Source: DI-hotos
	command := exec.CommandContext(ctx, nodeBinary, "-config", configPath, "-node", agentName)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("agent %s failed: %w", agentName, err)
	}
	return nil
}
