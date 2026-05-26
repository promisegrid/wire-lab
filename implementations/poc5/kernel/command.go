package kernel

import (
	"context"
	"flag"
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

// Command runs the poc5 local kernel.
func Command(args []string) error {
	flagSet := flag.NewFlagSet("poc5-kernel", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	appListen := flagSet.String("app-listen", "127.0.0.1:7201", "local app listener")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	if *nodeName == "" {
		return fmt.Errorf("--node is required")
	}
	localKernel := &Kernel{
		NodeName:           *nodeName,
		AppListen:          *appListen,
		ReceiveProtocolCID: ReceiveProtocolCID(),
		EvidenceLog:        lib.NewEvidenceLog(*nodeName, *nodeName+"-kernel"),
	}
	return localKernel.Run(context.Background())
}
