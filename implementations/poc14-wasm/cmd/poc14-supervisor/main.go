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
	"time"

	"promisegrid.dev/wire-lab/implementations/poc14-wasm/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.json", "POC14 config path")
	containerName := flag.String("container", "", "container name from config")
	kernelBinary := flag.String("kernel-bin", "poc14-kernel", "POC14 kernel binary path")
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
	return runContainerProcesses(ctx, cfg, *kernelBinary, *configPath, container)
}

func runContainerProcesses(ctx context.Context, cfg config.Config, kernelBinary, configPath string, container config.ContainerConfig) error {
	// Intent: The supervisor now proves inter-container promise delivery by
	// starting one container-local kernel plus independent app processes, rather
	// than hiding delivery inside one monolithic runtime object. Source:
	// DI-fumol; DI-sinur
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	kernelErrs := make(chan error, 1)
	go func() {
		kernelErrs <- runKernelProcess(ctx, kernelBinary, configPath, container.Name)
	}()
	select {
	case kernelErr := <-kernelErrs:
		if kernelErr != nil {
			return kernelErr
		}
		return fmt.Errorf("kernel for container %s exited before apps started", container.Name)
	case <-time.After(250 * time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}
	errs := make(chan error, len(container.Agents))
	var waitGroup sync.WaitGroup
	for _, agentName := range container.Agents {
		agentBinary, binaryErr := appBinaryForAgent(cfg, agentName)
		if binaryErr != nil {
			cancel()
			return binaryErr
		}
		waitGroup.Add(1)
		go func(localAgentName, localAgentBinary string) {
			defer waitGroup.Done()
			errs <- runAgentProcess(ctx, localAgentBinary, configPath, localAgentName)
		}(agentName, agentBinary)
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
	cancel()
	kernelErr := <-kernelErrs
	if firstErr == nil && kernelErr != nil {
		firstErr = kernelErr
	}
	return firstErr
}

func runKernelProcess(ctx context.Context, kernelBinary, configPath, containerName string) error {
	// Intent: Start exactly one local kernel per container. It is transport
	// plumbing for app promises, not a trust authority or workflow owner.
	// Source: DI-galin
	command := exec.CommandContext(ctx, kernelBinary, "-config", configPath, "-container", containerName)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		if ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("kernel %s failed: %w", containerName, err)
	}
	return nil
}

func runAgentProcess(ctx context.Context, agentBinary, configPath, agentName string) error {
	// Intent: A container supervisor starts independent local app processes
	// without sharing decision state; stdout/stderr pass through so the run log
	// remains auditable from Docker output. Source: DI-galin
	command := exec.CommandContext(ctx, agentBinary, "-config", configPath, "-node", agentName)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		if ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("agent %s failed: %w", agentName, err)
	}
	return nil
}

func appBinaryForAgent(cfg config.Config, agentName string) (string, error) {
	agent, ok := cfg.Agent(agentName)
	if !ok {
		return "", fmt.Errorf("unknown agent %q", agentName)
	}
	switch agent.Kind {
	case "":
		return "poc14-relationship-agent", nil
	case "fulfillment":
		return "poc14-fulfillment", nil
	case "postal_scale":
		return "poc14-postal-scale", nil
	case "ups_label_printer":
		return "poc14-ups-label-printer", nil
	case "printer_port":
		return "poc14-printer-port", nil
	case "accounting":
		return "poc14-accounting", nil
	case "wasm_agent":
		// Intent: POC14 adds heterogeneous runtime-adapter app roles as separate
		// supervised processes while preserving the same local-kernel routing
		// model used by POC12/POC13 agents. Source: DI-linof
		return "poc14-wasm-agent", nil
	case "stdio_agent":
		// Intent: The stdio adapter is the supervised app process for a worker
		// that sends and receives envelopes only over stdin/stdout. Source:
		// DI-linof
		return "poc14-stdio-adapter", nil
	default:
		return "", fmt.Errorf("agent %q has unsupported app kind %q", agentName, agent.Kind)
	}
}
