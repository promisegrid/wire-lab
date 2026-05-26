package hello

import (
	"context"
	"flag"
	"fmt"
)

// Main parses CLI flags and runs one hello app role.
func Main(ctx context.Context, args []string) error {
	flagSet := flag.NewFlagSet("poc4-hello", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	appName := flagSet.String("app", "", "local app name")
	kernelAddr := flagSet.String("kernel", "127.0.0.1:4101", "local kernel address")
	mode := flagSet.String("mode", "idle", "hello app mode: ask-signed or idle")
	targetNode := flagSet.String("target-node", "", "target node for ask-signed")
	targetApp := flagSet.String("target-app", "", "target signed app")
	text := flagSet.String("text", "hello from Alice", "hello text")
	if parseErr := flagSet.Parse(args); parseErr != nil {
		return parseErr
	}
	if *nodeName == "" {
		return fmt.Errorf("node is required")
	}
	if *mode == "ask-signed" && (*targetNode == "" || *targetApp == "") {
		return fmt.Errorf("target-node and target-app are required for ask-signed")
	}
	select {
	case <-ctx.Done():
		return nil
	default:
		return HelloApp{NodeName: *nodeName, AppName: *appName, KernelAddr: *kernelAddr, Mode: *mode, TargetNode: *targetNode, TargetApp: *targetApp, Text: *text}.Run()
	}
}
