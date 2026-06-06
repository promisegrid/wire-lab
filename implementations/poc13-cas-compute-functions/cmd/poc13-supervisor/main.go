package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

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
	var wg sync.WaitGroup
	errs := make(chan error, len(container.Agents))
	for _, agentName := range container.Agents {
		agent, agentOK := cfg.Agent(agentName)
		if !agentOK {
			return fmt.Errorf("unknown agent %s", agentName)
		}
		wg.Add(1)
		go func(localAgent poc13.AgentConfig) {
			defer wg.Done()
			runner := poc13.NewRunner(cfg, localAgent, poc13.LiveOrScriptedDecider{})
			if err := runner.Run(ctx); err != nil {
				errs <- err
			}
		}(agent)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
