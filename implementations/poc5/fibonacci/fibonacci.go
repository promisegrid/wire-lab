package fibonacci

import (
	"fmt"
	"strconv"

	"promisegrid.dev/wire-lab/implementations/poc5/kernel"
	"promisegrid.dev/wire-lab/implementations/poc5/lib"
	"promisegrid.dev/wire-lab/implementations/poc5/relay"
)

// FibonacciApp is an app agent that can promise to calculate or receive one
// fibonacci result.
type FibonacciApp struct {
	NodeName   string
	AppName    string
	KernelAddr string
	Mode       string
	TargetNode string
	TargetApp  string
	N          int
}

// Run executes one bounded fibonacci role.
func (fibonacciApp FibonacciApp) Run() error {
	if fibonacciApp.AppName == "" {
		fibonacciApp.AppName = fibonacciApp.NodeName + "-fibonacci-app"
	}
	switch fibonacciApp.Mode {
	case "client":
		return fibonacciApp.runClient()
	case "serve":
		return fibonacciApp.runServe()
	default:
		return fmt.Errorf("unknown fibonacci mode %q", fibonacciApp.Mode)
	}
}

func (fibonacciApp FibonacciApp) runClient() error {
	// Intent: Computation is modeled as reciprocal promises: Alice promises to
	// receive a result, then Carol may promise to calculate and return one.
	// Source: DI-rarim
	frameConn, registerErr := lib.RegisterReceiver(fibonacciApp.KernelAddr, kernel.ReceiveProtocolCID(), fibonacciApp.NodeName, fibonacciApp.AppName, ProtocolCID(), "I promise to receive one fibonacci_result_v1 and judge the computation locally.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	requestEnvelope, requestErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "fibonacci_request_v1",
		"from":      fibonacciApp.AppName,
		"from_node": fibonacciApp.NodeName,
		"to":        fibonacciApp.TargetApp,
		"n":         strconv.Itoa(fibonacciApp.N),
	})
	if requestErr != nil {
		return requestErr
	}
	requestBytes, sendErr := relay.SendInnerViaLocalRelay(fibonacciApp.KernelAddr, fibonacciApp.NodeName, fibonacciApp.AppName, fibonacciApp.TargetNode, fibonacciApp.TargetApp, requestEnvelope, "")
	if sendErr != nil {
		return sendErr
	}
	_, kind, fields, _, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	if kind != "fibonacci_result_v1" {
		return fmt.Errorf("fibonacci client got unexpected kind %s", kind)
	}
	expected := fibonacciNumber(fibonacciApp.N)
	result, parseErr := strconv.Atoi(fields["result"])
	if parseErr != nil {
		return parseErr
	}
	if fields["request_hash"] != lib.HashExactBytes(requestBytes) || result != expected {
		return fmt.Errorf("fibonacci result did not match local judgment")
	}
	fmt.Printf("%s judged fibonacci kept from %s: fib(%d)=%d\n", fibonacciApp.AppName, fields["from"], fibonacciApp.N, result)
	return nil
}

func (fibonacciApp FibonacciApp) runServe() error {
	frameConn, registerErr := lib.RegisterReceiver(fibonacciApp.KernelAddr, kernel.ReceiveProtocolCID(), fibonacciApp.NodeName, fibonacciApp.AppName, ProtocolCID(), "I promise to receive one fibonacci_request_v1 and return the result if the requester promises to receive it.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	envelope, kind, fields, requestBytes, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	if !envelope.ProtocolCID.Equal(ProtocolCID()) || kind != "fibonacci_request_v1" {
		return fmt.Errorf("fibonacci server got unexpected kind %s", kind)
	}
	n, parseErr := strconv.Atoi(fields["n"])
	if parseErr != nil {
		return parseErr
	}
	resultEnvelope, resultErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":         "fibonacci_result_v1",
		"from":         fibonacciApp.AppName,
		"from_node":    fibonacciApp.NodeName,
		"to":           fields["from"],
		"to_node":      fields["from_node"],
		"request_hash": lib.HashExactBytes(requestBytes),
		"n":            strconv.Itoa(n),
		"result":       strconv.Itoa(fibonacciNumber(n)),
	})
	if resultErr != nil {
		return resultErr
	}
	if _, sendErr := relay.SendInnerViaLocalRelay(fibonacciApp.KernelAddr, fibonacciApp.NodeName, fibonacciApp.AppName, fields["from_node"], fields["from"], resultEnvelope, lib.HashExactBytes(requestBytes)); sendErr != nil {
		return sendErr
	}
	fmt.Printf("%s judged fibonacci request from %s kept: fib(%d)=%d\n", fibonacciApp.AppName, fields["from"], n, fibonacciNumber(n))
	return nil
}

func fibonacciNumber(n int) int {
	if n < 2 {
		return n
	}
	previous := 0
	current := 1
	for index := 2; index <= n; index++ {
		next := previous + current
		previous = current
		current = next
	}
	return current
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
