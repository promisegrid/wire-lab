package kernel

import (
	"context"
	"flag"
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

// Command runs the poc3 kernel command.
func Command(args []string) error {
	flagSet := flag.NewFlagSet("poc3-kernel", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	appListen := flagSet.String("app-listen", "127.0.0.1:7101", "local app listener")
	peerListen := flagSet.String("peer-listen", "0.0.0.0:9100", "peer kernel listener")
	peerAddress := flagSet.String("peer", "", "remote peer kernel address")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	if *nodeName == "" {
		return fmt.Errorf("--node is required")
	}
	kernel := &Kernel{
		NodeName:           *nodeName,
		AppListen:          *appListen,
		PeerListen:         *peerListen,
		PeerAddress:        *peerAddress,
		ReceiveProtocolCID: ReceiveProtocolCID(),
		EvidenceLog:        lib.NewEvidenceLog(*nodeName),
	}
	return kernel.Run(context.Background())
}
