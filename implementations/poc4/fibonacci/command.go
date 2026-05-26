package fibonacci

import (
	"context"
	"flag"
	"fmt"
)

// Main parses CLI flags and runs one fibonacci app role.
func Main(ctx context.Context, args []string) error {
	flagSet := flag.NewFlagSet("poc4-fibonacci", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	appName := flagSet.String("app", "", "local app name")
	kernelAddr := flagSet.String("kernel", "127.0.0.1:4101", "local kernel address")
	mode := flagSet.String("mode", "serve", "fibonacci mode: client or serve")
	targetNode := flagSet.String("target-node", "", "target node for client mode")
	targetApp := flagSet.String("target-app", "", "target app for client mode")
	n := flagSet.Int("n", 10, "fibonacci input")
	if parseErr := flagSet.Parse(args); parseErr != nil {
		return parseErr
	}
	if *nodeName == "" {
		return fmt.Errorf("node is required")
	}
	if *mode == "client" && (*targetNode == "" || *targetApp == "") {
		return fmt.Errorf("target-node and target-app are required for client mode")
	}
	select {
	case <-ctx.Done():
		return nil
	default:
		return FibonacciApp{NodeName: *nodeName, AppName: *appName, KernelAddr: *kernelAddr, Mode: *mode, TargetNode: *targetNode, TargetApp: *targetApp, N: *n}.Run()
	}
}
