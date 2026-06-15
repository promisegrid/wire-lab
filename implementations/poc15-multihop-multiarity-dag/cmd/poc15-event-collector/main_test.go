package main

import "testing"

func TestCollectorListenAddressUsesConfiguredPort(t *testing.T) {
	listenAddress, err := collectorListenAddress("event-collector:9200")
	if err != nil {
		t.Fatalf("collector listen address: %v", err)
	}
	if listenAddress != ":9200" {
		t.Fatalf("listen address = %q, want :9200", listenAddress)
	}
}
