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
	nodeName := flag.String("node", "", "printer port app name")
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
	if agent.Kind != "printer_port" {
		return fmt.Errorf("node %q is kind %q, want printer_port", agent.Name, agent.Kind)
	}
	// Intent: printer_port is a local kernel-role app around a simulated printer
	// device. It promises capability-token issue and bounded local print
	// event records, while the message kernel only routes exact pCID frames.
	// Source: DI-pohaj
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
