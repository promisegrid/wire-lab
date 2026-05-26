package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := runMain(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func runMain(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing argv")
	}
	commandName := filepath.Base(args[0])
	remainingArgs := args[1:]
	if commandName != "poc2-kernel" && commandName != "poc2-hello" && len(remainingArgs) > 0 {
		candidate := remainingArgs[0]
		if candidate == "kernel" || candidate == "hello" {
			commandName = "poc2-" + candidate
			remainingArgs = remainingArgs[1:]
		}
	}
	protocolCID, protocolErr := HelloProtocolCID()
	if protocolErr != nil {
		return protocolErr
	}
	switch {
	case strings.HasSuffix(commandName, "poc2-kernel"):
		return runKernelCommand(remainingArgs, protocolCID)
	case strings.HasSuffix(commandName, "poc2-hello"):
		return runHelloCommand(remainingArgs, protocolCID)
	default:
		return fmt.Errorf("invoke as poc2-kernel, poc2-hello, or use subcommand kernel|hello")
	}
}

func runKernelCommand(args []string, protocolCID ProtocolCID) error {
	flagSet := flag.NewFlagSet("poc2-kernel", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	appListen := flagSet.String("app-listen", "127.0.0.1:7001", "local app listener")
	peerListen := flagSet.String("peer-listen", "0.0.0.0:9000", "peer kernel listener")
	peerAddress := flagSet.String("peer", "", "remote peer kernel address")
	evidencePath := flagSet.String("evidence", "", "optional JSONL evidence path")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	if *nodeName == "" {
		return fmt.Errorf("--node is required")
	}
	evidenceLog, cleanup, logErr := NewEvidenceLog(*nodeName, *evidencePath)
	if logErr != nil {
		return logErr
	}
	defer func() {
		if cleanupErr := cleanup(); cleanupErr != nil {
			fmt.Fprintln(os.Stderr, cleanupErr.Error())
		}
	}()
	kernel := &Kernel{
		NodeName:    *nodeName,
		AppListen:   *appListen,
		PeerListen:  *peerListen,
		PeerAddress: *peerAddress,
		ProtocolCID: protocolCID,
		EvidenceLog: evidenceLog,
	}
	return kernel.Run(context.Background())
}

func runHelloCommand(args []string, protocolCID ProtocolCID) error {
	flagSet := flag.NewFlagSet("poc2-hello", flag.ContinueOnError)
	nodeName := flagSet.String("node", "", "local node name")
	kernelAddress := flagSet.String("kernel", "127.0.0.1:7001", "local kernel app address")
	mode := flagSet.String("mode", "", "send or receive")
	destination := flagSet.String("to", "", "destination node for send mode")
	text := flagSet.String("text", "hello", "hello text")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	if *nodeName == "" {
		return fmt.Errorf("--node is required")
	}
	helloApp := HelloApp{
		NodeName:    *nodeName,
		KernelAddr:  *kernelAddress,
		Mode:        *mode,
		Destination: *destination,
		Text:        *text,
		ProtocolCID: protocolCID,
	}
	return helloApp.Run()
}
