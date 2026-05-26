package main

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"
)

func TestTwoKernelHelloFlow(t *testing.T) {
	protocolCID, cidErr := HelloProtocolCID()
	if cidErr != nil {
		t.Fatalf("HelloProtocolCID: %v", cidErr)
	}
	aliceAppAddress := reserveTCPAddress(t)
	alicePeerAddress := reserveTCPAddress(t)
	bobAppAddress := reserveTCPAddress(t)
	bobPeerAddress := reserveTCPAddress(t)
	aliceLog, aliceCleanup, aliceLogErr := NewEvidenceLog("alice", "")
	if aliceLogErr != nil {
		t.Fatalf("NewEvidenceLog alice: %v", aliceLogErr)
	}
	defer func() {
		if cleanupErr := aliceCleanup(); cleanupErr != nil {
			t.Fatalf("alice cleanup: %v", cleanupErr)
		}
	}()
	bobLog, bobCleanup, bobLogErr := NewEvidenceLog("bob", "")
	if bobLogErr != nil {
		t.Fatalf("NewEvidenceLog bob: %v", bobLogErr)
	}
	defer func() {
		if cleanupErr := bobCleanup(); cleanupErr != nil {
			t.Fatalf("bob cleanup: %v", cleanupErr)
		}
	}()
	aliceKernel := &Kernel{NodeName: "alice", AppListen: aliceAppAddress, PeerListen: alicePeerAddress, PeerAddress: bobPeerAddress, ProtocolCID: protocolCID, EvidenceLog: aliceLog}
	bobKernel := &Kernel{NodeName: "bob", AppListen: bobAppAddress, PeerListen: bobPeerAddress, ProtocolCID: protocolCID, EvidenceLog: bobLog}
	testContext, cancel := context.WithCancel(context.Background())
	defer cancel()
	errorChannel := make(chan error, 2)
	go func() { errorChannel <- aliceKernel.Run(testContext) }()
	go func() { errorChannel <- bobKernel.Run(testContext) }()

	receiveDone := make(chan error, 1)
	go func() {
		receiver := HelloApp{NodeName: "bob", KernelAddr: bobAppAddress, Mode: "receive", Destination: "bob", Text: "receive", ProtocolCID: protocolCID}
		receiveDone <- receiver.Run()
	}()
	sender := HelloApp{NodeName: "alice", KernelAddr: aliceAppAddress, Mode: "send", Destination: "bob", Text: "hello from Alice", ProtocolCID: protocolCID}
	if sendErr := sender.Run(); sendErr != nil {
		t.Fatalf("sender run: %v", sendErr)
	}
	select {
	case receiveErr := <-receiveDone:
		if receiveErr != nil {
			t.Fatalf("receiver run: %v", receiveErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("receiver did not complete")
	}
	cancel()
	for completedRuns := 0; completedRuns < 2; completedRuns++ {
		select {
		case runErr := <-errorChannel:
			if runErr != nil {
				t.Fatalf("kernel run: %v", runErr)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("kernel did not stop")
		}
	}
}

func reserveTCPAddress(t *testing.T) string {
	listener, listenErr := net.Listen("tcp", "127.0.0.1:0")
	if listenErr != nil {
		if strings.Contains(listenErr.Error(), "operation not permitted") {
			t.Skipf("local TCP listeners are blocked in this sandbox: %v", listenErr)
		}
		t.Fatalf("reserve TCP address: %v", listenErr)
	}
	address := listener.Addr().String()
	if closeErr := listener.Close(); closeErr != nil {
		t.Fatalf("close reserved listener: %v", closeErr)
	}
	return address
}
