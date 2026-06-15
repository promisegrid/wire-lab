package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/config"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/decision"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/runtime"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.json", "POC15 config path")
	nodeName := flag.String("node", "", "accounting app name")
	flag.Parse()
	if *nodeName == "" {
		return fmt.Errorf("-node is required")
	}
	cfg, loadErr := config.Load(*configPath)
	if loadErr != nil {
		return loadErr
	}
	agent, ok := cfg.Agent(*nodeName)
	if !ok {
		return fmt.Errorf("unknown node %q", *nodeName)
	}
	if agent.Kind != "accounting" {
		return fmt.Errorf("node %q is kind %q, want accounting", agent.Name, agent.Kind)
	}
	// Intent: Accounting is a deterministic app process that promises address
	// and shipment-update event records; it is not a kernel-owned service registry.
	// Source: DI-galin
	liveClient := decision.NewLiveClient(
		cfg.ProviderBaseURL,
		cfg.APIKeyEnv,
		cfg.AgentModel,
		cfg.MonitorModel,
		cfg.ReasoningEffort,
		cfg.ServiceTier,
		cfg.RequestTimeout(),
	)
	node := runtime.NewNode(cfg, agent, liveClient, liveClient)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return node.Run(ctx)
}
