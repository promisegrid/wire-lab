package main

import (
	"flag"
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/config"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/sim"
)

func main() {
	configPath := flag.String("config", "config.json", "path to POC17 simulator config")
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}
	if err := sim.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "run simulator: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(cfg.RunDir())
}
