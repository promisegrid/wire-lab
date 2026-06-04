package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/config"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/decision"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/runtime"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.json", "POC12 config path")
	nodeName := flag.String("node", "", "agent node name")
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
	// Intent: Live runtime uses the config-selected provider and API-key
	// environment variable, while tests instantiate FakeDecider directly.
	// Source: DI-timah
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
