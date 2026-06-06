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
	agentName := flag.String("agent", "", "agent name from config")
	flag.Parse()
	if err := run(context.Background(), *configPath, *agentName); err != nil {
		fmt.Fprintf(os.Stderr, "poc13-agent: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, configPath, agentName string) error {
	if agentName == "" {
		return fmt.Errorf("agent is required")
	}
	cfg, loadErr := poc13.LoadConfig(configPath)
	if loadErr != nil {
		return loadErr
	}
	agent, ok := cfg.Agent(agentName)
	if !ok {
		return fmt.Errorf("unknown agent %s", agentName)
	}
	return poc13.NewRunner(cfg, agent, poc13.LiveOrScriptedDecider{}).Run(ctx)
}
