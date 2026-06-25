package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/decision"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/lifecycle"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/runtime"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.json", "POC16 config path")
	nodeName := flag.String("node", "", "UPS label printer app name")
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
	if agent.Kind != "ups_label_printer" {
		return fmt.Errorf("node %q is kind %q, want ups_label_printer", agent.Name, agent.Kind)
	}
	// Intent: The label printer remains an app process that promises label,
	// cost, and tracking event records; the kernel only delivers its pCID frames.
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
	return lifecycle.RunManaged(ctx, appLifecycleOptions(cfg, agent.Name), node.Run)
}

func appLifecycleOptions(cfg config.Config, agentName string) lifecycle.RoleOptions {
	containerName, _ := cfg.ContainerForAgent(agentName)
	return lifecycle.OptionsFromEnv("supervisor:"+containerName, "agent:"+agentName, lifecycle.RoleKindApp, cfg.RunID, pcid.NewRegistry().MustCID(pcid.LocalLifecycleV1), cfg.ShutdownGrace(), true, os.Stdin)
}
