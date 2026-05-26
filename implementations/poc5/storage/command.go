package storage

import (
	"context"
	"flag"
	"fmt"
)

// Main parses CLI flags and runs one storage app role.
func Main(ctx context.Context, args []string) error {
	flagSet := flag.NewFlagSet("poc5-storage", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	appName := flagSet.String("app", "", "local app name")
	kernelAddr := flagSet.String("kernel", "127.0.0.1:4101", "local kernel address")
	mode := flagSet.String("mode", "serve", "storage mode: client, trust-client, serve, or serve-break")
	targetNode := flagSet.String("target-node", "", "target node for client mode")
	targetApp := flagSet.String("target-app", "", "target app for client mode")
	fallbackNode := flagSet.String("fallback-node", "", "fallback node for trust-client mode")
	fallbackApp := flagSet.String("fallback-app", "", "fallback app for trust-client mode")
	key := flagSet.String("key", "poc5-key", "storage key")
	value := flagSet.String("value", "poc5-value", "storage value")
	if parseErr := flagSet.Parse(args); parseErr != nil {
		return parseErr
	}
	if *nodeName == "" {
		return fmt.Errorf("node is required")
	}
	if *mode == "client" && (*targetNode == "" || *targetApp == "") {
		return fmt.Errorf("target-node and target-app are required for client mode")
	}
	if *mode == "trust-client" && (*targetNode == "" || *targetApp == "" || *fallbackNode == "" || *fallbackApp == "") {
		return fmt.Errorf("target-node, target-app, fallback-node, and fallback-app are required for trust-client mode")
	}
	select {
	case <-ctx.Done():
		return nil
	default:
		return StorageApp{NodeName: *nodeName, AppName: *appName, KernelAddr: *kernelAddr, Mode: *mode, TargetNode: *targetNode, TargetApp: *targetApp, FallbackNode: *fallbackNode, FallbackApp: *fallbackApp, Key: *key, Value: *value}.Run()
	}
}
