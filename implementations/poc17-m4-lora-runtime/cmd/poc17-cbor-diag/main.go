package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/protocol"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: poc17-cbor-diag <message.cbor>")
		os.Exit(2)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read CBOR: %v\n", err)
		os.Exit(1)
	}
	msg, err := protocol.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse CBOR: %v\n", err)
		os.Exit(1)
	}
	protocolName := msg.ProtocolName
	if protocolName == "" {
		protocolName = "unknown"
	}
	fmt.Printf("grid([42(%s /* %s */), h'%s'", msg.PCID, protocolName, hex.EncodeToString(msg.Payload))
	if msg.Proof != nil {
		fmt.Printf(", h'%s'", hex.EncodeToString(msg.Proof))
	}
	fmt.Println("])")
}
