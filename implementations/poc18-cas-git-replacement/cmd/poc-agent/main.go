// Command poc-agent runs one deterministic POC18 Docker agent.
//
// Intent: POC18 remediation needs each named agent to own a private CAS,
// listener, peer table, and observer stream inside its own container instead of
// sharing one in-process fixture. Source: DI-koriz
package main

import (
	"flag"
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/agent"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "poc-agent: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("poc-agent", flag.ContinueOnError)
	name := flags.String("agent", "", "agent name")
	listen := flags.String("listen", "", "agent TCP listen address")
	casRoot := flags.String("cas", "", "agent-local CAS root")
	collector := flags.String("collector", "", "observer collector address")
	var peers agent.PeerFlags
	flags.Var(&peers, "peer", "peer in name=address form; may repeat")
	if err := flags.Parse(args); err != nil {
		return err
	}
	peerMap, peerErr := agent.ParsePeers(peers)
	if peerErr != nil {
		return peerErr
	}
	return agent.Run(agent.Config{Name: *name, Listen: *listen, CASRoot: *casRoot, Collector: *collector, Peers: peerMap})
}
