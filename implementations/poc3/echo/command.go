package echo

import (
	"flag"
	"fmt"
)

// Command runs the poc3 echo app command.
func Command(args []string) error {
	flagSet := flag.NewFlagSet("poc3-echo", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	kernelAddress := flagSet.String("kernel", "127.0.0.1:7101", "local kernel app address")
	mode := flagSet.String("mode", "", "ask or serve")
	destination := flagSet.String("to", "", "destination node for ask mode")
	text := flagSet.String("text", "echo", "echo text")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	if *nodeName == "" {
		return fmt.Errorf("--node is required")
	}
	return EchoApp{NodeName: *nodeName, KernelAddr: *kernelAddress, Mode: *mode, Destination: *destination, Text: *text}.Run()
}
