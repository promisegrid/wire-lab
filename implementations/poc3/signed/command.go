package signed

import (
	"flag"
	"fmt"
)

// Command runs the poc3 signed app command.
func Command(args []string) error {
	flagSet := flag.NewFlagSet("poc3-signed", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	kernelAddress := flagSet.String("kernel", "127.0.0.1:7101", "local kernel app address")
	mode := flagSet.String("mode", "", "send or receive")
	destination := flagSet.String("to", "", "destination node for send mode")
	text := flagSet.String("text", "signed hello", "signed message text")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	if *nodeName == "" {
		return fmt.Errorf("--node is required")
	}
	return SignedApp{NodeName: *nodeName, KernelAddr: *kernelAddress, Mode: *mode, Destination: *destination, Text: *text}.Run()
}
