package relay

import (
	"context"
	"flag"
	"fmt"
)

// Command runs the poc4 relay app.
func Command(args []string) error {
	flagSet := flag.NewFlagSet("poc4-relay", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	kernelAddress := flagSet.String("kernel", "127.0.0.1:7201", "local kernel app address")
	listenAddress := flagSet.String("listen", "0.0.0.0:9200", "relay neighbor listener")
	routeText := flagSet.String("routes", "", "local route promises target=addr,target=addr")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	if *nodeName == "" {
		return fmt.Errorf("--node is required")
	}
	routes, routeErr := ParseRoutes(*routeText)
	if routeErr != nil {
		return routeErr
	}
	return RelayApp{NodeName: *nodeName, KernelAddr: *kernelAddress, ListenAddr: *listenAddress, RouteTable: routes}.Run(context.Background())
}
