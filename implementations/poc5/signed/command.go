package signed

import (
	"context"
	"flag"
	"fmt"
)

// Main parses CLI flags and runs one signed app role.
func Main(ctx context.Context, args []string) error {
	flagSet := flag.NewFlagSet("poc5-signed", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	appName := flagSet.String("app", "", "local app name")
	kernelAddr := flagSet.String("kernel", "127.0.0.1:4101", "local kernel address")
	mode := flagSet.String("mode", "serve", "signed app mode: serve or idle")
	if parseErr := flagSet.Parse(args); parseErr != nil {
		return parseErr
	}
	if *nodeName == "" {
		return fmt.Errorf("node is required")
	}
	select {
	case <-ctx.Done():
		return nil
	default:
		return SignedApp{NodeName: *nodeName, AppName: *appName, KernelAddr: *kernelAddr, Mode: *mode}.Run()
	}
}
