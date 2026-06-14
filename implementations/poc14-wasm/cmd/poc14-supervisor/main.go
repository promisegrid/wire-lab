package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc14-wasm/config"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/decision"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/eventstream"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.json", "POC14 config path")
	containerName := flag.String("container", "", "container name from config")
	kernelBinary := flag.String("kernel-bin", "poc14-kernel", "POC14 kernel binary path")
	flag.Parse()
	if *containerName == "" {
		return fmt.Errorf("-container is required")
	}
	cfg, loadErr := config.Load(*configPath)
	if loadErr != nil {
		return loadErr
	}
	container, containerFound := cfg.Container(*containerName)
	if !containerFound {
		return fmt.Errorf("unknown container %q", *containerName)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return runContainerProcesses(ctx, cfg, *kernelBinary, *configPath, container)
}

func runContainerProcesses(ctx context.Context, cfg config.Config, kernelBinary, configPath string, container config.ContainerConfig) error {
	// Intent: The supervisor now proves inter-container promise delivery by
	// starting one container-local kernel plus independent app processes, rather
	// than hiding delivery inside one monolithic runtime object. Source:
	// DI-fumol; DI-sinur
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	collectorClient, collectorErr := eventstream.Dial(ctx, cfg.EventCollectorAddress, cfg.StartupDelay()+30*time.Second)
	if collectorErr != nil {
		return collectorErr
	}
	defer closeCollectorClient(collectorClient)
	forwarder := newProcessForwarder(container.Name, collectorClient)
	defer forwarder.sendDone("container supervisor exited")
	kernelErrs := make(chan error, 1)
	go func() {
		kernelErrs <- runKernelProcess(ctx, forwarder, kernelBinary, configPath, container.Name)
	}()
	select {
	case kernelErr := <-kernelErrs:
		if kernelErr != nil {
			return kernelErr
		}
		return fmt.Errorf("kernel for container %s exited before apps started", container.Name)
	case <-time.After(250 * time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}
	errs := make(chan error, len(container.Agents))
	var waitGroup sync.WaitGroup
	for _, agentName := range container.Agents {
		agentBinary, binaryErr := appBinaryForAgent(cfg, agentName)
		if binaryErr != nil {
			cancel()
			return binaryErr
		}
		waitGroup.Add(1)
		go func(localAgentName, localAgentBinary string) {
			defer waitGroup.Done()
			errs <- runAgentProcess(ctx, forwarder, localAgentBinary, configPath, localAgentName)
		}(agentName, agentBinary)
	}
	go func() {
		waitGroup.Wait()
		close(errs)
	}()
	var firstErr error
	for agentErr := range errs {
		if agentErr != nil && firstErr == nil {
			firstErr = agentErr
			cancel()
		}
	}
	cancel()
	kernelErr := <-kernelErrs
	if firstErr == nil && kernelErr != nil {
		firstErr = kernelErr
	}
	if forwardErr := forwarder.err(); firstErr == nil && forwardErr != nil {
		firstErr = forwardErr
	}
	return firstErr
}

func runKernelProcess(ctx context.Context, forwarder *processForwarder, kernelBinary, configPath, containerName string) error {
	// Intent: Start exactly one local kernel per container. It is transport
	// plumbing for app promises, not a trust authority or workflow owner.
	// Source: DI-galin
	return runProcess(ctx, forwarder, "kernel:"+containerName, kernelBinary, "-config", configPath, "-container", containerName)
}

func runAgentProcess(ctx context.Context, forwarder *processForwarder, agentBinary, configPath, agentName string) error {
	// Intent: A container supervisor starts independent local app processes
	// without sharing decision state; stdout/stderr pass through so the run log
	// remains auditable from Docker output. Source: DI-galin
	return runProcess(ctx, forwarder, "agent:"+agentName, agentBinary, "-config", configPath, "-node", agentName)
}

func runProcess(ctx context.Context, forwarder *processForwarder, roleName, binary string, args ...string) error {
	// Intent: Child processes still write ordinary stdout/stderr, but the
	// supervisor copies JSON event records to the observer-only collector instead
	// of relying on a Docker volume visible to agents. Source: DI-dirat
	command := exec.CommandContext(ctx, binary, args...)
	stdout, stdoutErr := command.StdoutPipe()
	if stdoutErr != nil {
		return stdoutErr
	}
	stderr, stderrErr := command.StderrPipe()
	if stderrErr != nil {
		return stderrErr
	}
	if startErr := command.Start(); startErr != nil {
		return startErr
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go forwarder.copyOutput(&waitGroup, roleName, "stdout", stdout, os.Stdout)
	go forwarder.copyOutput(&waitGroup, roleName, "stderr", stderr, os.Stderr)
	waitErr := command.Wait()
	waitGroup.Wait()
	if waitErr != nil {
		if ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("%s failed: %w", roleName, waitErr)
	}
	return nil
}

type processForwarder struct {
	containerName string
	client        *eventstream.Client

	mu       sync.Mutex
	firstErr error
	doneSent bool
}

func newProcessForwarder(containerName string, client *eventstream.Client) *processForwarder {
	return &processForwarder{containerName: containerName, client: client}
}

func (forwarder *processForwarder) copyOutput(waitGroup *sync.WaitGroup, roleName, streamName string, reader io.Reader, writer io.Writer) {
	defer waitGroup.Done()
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if _, writeErr := fmt.Fprintln(writer, line); writeErr != nil {
			forwarder.rememberErr(writeErr)
		}
		forwarder.forwardEventLine(roleName, streamName, line)
	}
	if scanErr := scanner.Err(); scanErr != nil {
		// Intent: When the supervisor cancels the local kernel after app
		// processes exit, Go may report closed stdout/stderr pipes. That is normal
		// process-lifecycle cleanup, not a broken promise or failed container.
		// Source: DI-dirat
		if !errors.Is(scanErr, os.ErrClosed) {
			forwarder.rememberErr(scanErr)
		}
	}
}

func (forwarder *processForwarder) forwardEventLine(roleName, streamName, line string) {
	var event decision.Event
	if unmarshalErr := json.Unmarshal([]byte(line), &event); unmarshalErr != nil {
		return
	}
	if event.Observer == "" || event.Event == "" {
		return
	}
	sendErr := forwarder.client.Send(eventstream.Record{
		Kind:   eventstream.KindEvent,
		Source: forwarder.containerName + "/" + roleName + "/" + streamName,
		Event:  &event,
	})
	if sendErr != nil {
		forwarder.rememberErr(sendErr)
	}
}

func (forwarder *processForwarder) sendDone(detail string) {
	forwarder.mu.Lock()
	if forwarder.doneSent {
		forwarder.mu.Unlock()
		return
	}
	forwarder.doneSent = true
	forwarder.mu.Unlock()
	sendErr := forwarder.client.Send(eventstream.Record{
		Kind:   eventstream.KindSupervisorDone,
		Source: forwarder.containerName,
		Detail: detail,
	})
	if sendErr != nil {
		forwarder.rememberErr(sendErr)
	}
}

func (forwarder *processForwarder) rememberErr(err error) {
	forwarder.mu.Lock()
	defer forwarder.mu.Unlock()
	if forwarder.firstErr == nil {
		forwarder.firstErr = err
	}
}

func (forwarder *processForwarder) err() error {
	forwarder.mu.Lock()
	defer forwarder.mu.Unlock()
	return forwarder.firstErr
}

func closeCollectorClient(client *eventstream.Client) {
	if closeErr := client.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "close event collector client: %v\n", closeErr)
	}
}

func appBinaryForAgent(cfg config.Config, agentName string) (string, error) {
	agent, ok := cfg.Agent(agentName)
	if !ok {
		return "", fmt.Errorf("unknown agent %q", agentName)
	}
	switch agent.Kind {
	case "":
		return "poc14-relationship-agent", nil
	case "fulfillment":
		return "poc14-fulfillment", nil
	case "postal_scale":
		return "poc14-postal-scale", nil
	case "ups_label_printer":
		return "poc14-ups-label-printer", nil
	case "printer_port":
		return "poc14-printer-port", nil
	case "accounting":
		return "poc14-accounting", nil
	case "wasm_agent":
		// Intent: POC14 adds heterogeneous runtime-adapter app roles as separate
		// supervised processes while preserving the same local-kernel routing
		// model used by POC12/POC13 agents. Source: DI-linof
		return "poc14-wasm-agent", nil
	case "stdio_agent":
		// Intent: The stdio adapter is the supervised app process for a worker
		// that sends and receives envelopes only over stdin/stdout. Source:
		// DI-linof
		return "poc14-stdio-adapter", nil
	default:
		return "", fmt.Errorf("agent %q has unsupported app kind %q", agentName, agent.Kind)
	}
}
