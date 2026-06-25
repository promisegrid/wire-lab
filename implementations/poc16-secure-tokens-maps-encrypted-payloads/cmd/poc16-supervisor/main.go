package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/decision"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/eventstream"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/lifecycle"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.json", "POC16 config path")
	containerName := flag.String("container", "", "container name from config")
	kernelBinary := flag.String("kernel-bin", "poc16-kernel", "POC16 kernel binary path")
	parserRoleBinary := flag.String("parser-role-bin", "poc16-parser-role", "POC16 parser-role binary path")
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
	return runContainerProcesses(ctx, cfg, *kernelBinary, *parserRoleBinary, *configPath, container)
}

func runContainerProcesses(ctx context.Context, cfg config.Config, kernelBinary, parserRoleBinary, configPath string, container config.ContainerConfig) error {
	// Intent: The supervisor now proves inter-container promise delivery by
	// starting one container-local kernel plus independent app processes, rather
	// than hiding delivery inside one monolithic runtime object. Source:
	// DI-fumol; DI-sinur
	lifecycleCtx, cancelAll := context.WithCancel(ctx)
	defer cancelAll()
	kernelCtx, cancelKernel := context.WithCancel(lifecycleCtx)
	defer cancelKernel()
	parserCtx, cancelParser := context.WithCancel(lifecycleCtx)
	defer cancelParser()
	agentCtx, cancelAgents := context.WithCancel(lifecycleCtx)
	defer cancelAgents()
	// Intent: Use the same configured shutdown grace for supervisor SIGTERM
	// escalation that app runtimes use for local receive-promise grace, so verbose
	// parser/kernel terminal records are not truncated by a stale hard-coded kill
	// window. Source: DI-titik
	shutdownGrace := cfg.ShutdownGrace()
	collectorClient, collectorErr := eventstream.Dial(lifecycleCtx, cfg.EventCollectorAddress, cfg.StartupDelay()+30*time.Second)
	if collectorErr != nil {
		return collectorErr
	}
	defer closeCollectorClient(collectorClient)
	forwarder := newProcessForwarder(container.Name, collectorClient)
	defer forwarder.sendDone("container supervisor exited")
	lifecycleSupervisor, lifecycleErr := lifecycle.NewSupervisor("supervisor:"+container.Name, cfg.RunID, pcid.NewRegistry().MustCID(pcid.LocalLifecycleV1), shutdownGrace, func(event decision.Event) {
		forwarder.sendEvent("supervisor/lifecycle", event)
	})
	if lifecycleErr != nil {
		return lifecycleErr
	}
	defer closeLifecycleSupervisor(lifecycleSupervisor)
	kernelErrs := make(chan error, 1)
	kernelRoleID := "kernel:" + container.Name
	go func() {
		kernelErrs <- runKernelProcess(kernelCtx, forwarder, shutdownGrace, kernelBinary, configPath, container.Name, lifecycleSupervisor.EnvFor(protocol.LifecycleChannelTCP), nil)
	}()
	select {
	case kernelErr := <-kernelErrs:
		if kernelErr != nil {
			return kernelErr
		}
		return fmt.Errorf("kernel for container %s exited before apps started", container.Name)
	case <-time.After(250 * time.Millisecond):
	case <-lifecycleCtx.Done():
		return lifecycleCtx.Err()
	}
	parserErrs := make(chan error, 1)
	parserRoleID := "parser:" + container.Name
	go func() {
		parserErrs <- runParserRoleProcess(parserCtx, forwarder, shutdownGrace, parserRoleBinary, configPath, container.Name, lifecycleSupervisor.EnvFor(protocol.LifecycleChannelParserPath), nil)
	}()
	select {
	case parserErr := <-parserErrs:
		cancelAll()
		if parserErr != nil {
			<-kernelErrs
			return parserErr
		}
		<-kernelErrs
		return fmt.Errorf("parser role for container %s exited before apps started", container.Name)
	case kernelErr := <-kernelErrs:
		cancelAll()
		if kernelErr != nil {
			<-parserErrs
			return kernelErr
		}
		<-parserErrs
		return fmt.Errorf("kernel for container %s exited before parser role started apps", container.Name)
	case <-time.After(250 * time.Millisecond):
	case <-lifecycleCtx.Done():
		cancelAll()
		<-kernelErrs
		<-parserErrs
		return lifecycleCtx.Err()
	}
	errs := make(chan error, len(container.Agents))
	var waitGroup sync.WaitGroup
	appRoleIDs := make([]string, 0, len(container.Agents))
	for _, agentName := range container.Agents {
		agentBinary, binaryErr := appBinaryForAgent(cfg, agentName)
		if binaryErr != nil {
			cancelAll()
			<-parserErrs
			<-kernelErrs
			return binaryErr
		}
		agent, agentFound := cfg.Agent(agentName)
		if !agentFound {
			cancelAll()
			<-parserErrs
			<-kernelErrs
			return fmt.Errorf("unknown agent %q", agentName)
		}
		roleID := "agent:" + agentName
		appRoleIDs = append(appRoleIDs, roleID)
		channelProfile := lifecycleProfileForAgent(agent)
		waitGroup.Add(1)
		go func(localAgentName, localAgentBinary, localRoleID, localChannelProfile string) {
			defer waitGroup.Done()
			errs <- runAgentProcess(agentCtx, forwarder, shutdownGrace, localAgentBinary, configPath, localAgentName, lifecycleSupervisor.EnvFor(localChannelProfile), func(stdin io.WriteCloser) {
				if localChannelProfile == protocol.LifecycleChannelStdio {
					lifecycleSupervisor.RegisterStdin(localRoleID, stdin)
				}
			})
		}(agentName, agentBinary, roleID, channelProfile)
	}
	go func() {
		waitGroup.Wait()
		close(errs)
	}()
	readyCtx, readyCancel := context.WithTimeout(lifecycleCtx, cfg.MonitorWaitTimeout())
	defer readyCancel()
	if readyErr := lifecycleSupervisor.WaitReady(readyCtx, append([]string{kernelRoleID, parserRoleID}, appRoleIDs...)); readyErr != nil {
		forwarder.sendEvent("supervisor/lifecycle", decision.Event{Observer: "supervisor:" + container.Name, Event: "local_lifecycle_sigterm_fallback_used", Outcome: "broken", Detail: readyErr.Error()})
		cancelAll()
		<-parserErrs
		<-kernelErrs
		return readyErr
	}
	appInvokeCtx, appInvokeCancel := context.WithTimeout(lifecycleCtx, shutdownGrace)
	defer appInvokeCancel()
	for _, roleID := range appRoleIDs {
		if invokeErr := lifecycleSupervisor.Invoke(appInvokeCtx, roleID, "run_complete", ""); invokeErr != nil {
			forwarder.sendEvent("supervisor/lifecycle", decision.Event{Observer: "supervisor:" + container.Name, Event: "local_lifecycle_sigterm_fallback_used", Outcome: "broken", Peer: roleID, Detail: invokeErr.Error()})
			cancelAgents()
			break
		}
	}
	if fulfillErr := lifecycleSupervisor.WaitFulfilled(appInvokeCtx, appRoleIDs); fulfillErr != nil {
		forwarder.sendEvent("supervisor/lifecycle", decision.Event{Observer: "supervisor:" + container.Name, Event: "local_lifecycle_sigterm_fallback_used", Outcome: "broken", Detail: fulfillErr.Error()})
		cancelAgents()
	}
	var firstErr error
	for agentErr := range errs {
		if agentErr != nil && firstErr == nil {
			firstErr = agentErr
			cancelAgents()
		}
	}
	cancelAgents()
	// Intent: Parser roles own app-facing payload parsing and the parser/kernel
	// control stream. Shutting the parser down before the transport kernel lets
	// it close those local streams and emit terminal records before the kernel
	// closes peer transport sessions. Source: DI-vazoz
	parserAddress, parserAddressFound := parserRoleAddress(cfg, container.Name)
	parserInvokeCtx, parserInvokeCancel := context.WithTimeout(lifecycleCtx, shutdownGrace)
	defer parserInvokeCancel()
	if !parserAddressFound {
		if firstErr == nil {
			firstErr = fmt.Errorf("no parser lifecycle address for container %s", container.Name)
		}
		cancelParser()
	} else if invokeErr := lifecycleSupervisor.Invoke(parserInvokeCtx, parserRoleID, "apps_complete", parserAddress); invokeErr != nil {
		forwarder.sendEvent("supervisor/lifecycle", decision.Event{Observer: "supervisor:" + container.Name, Event: "local_lifecycle_sigterm_fallback_used", Outcome: "broken", Peer: parserRoleID, Detail: invokeErr.Error()})
		if firstErr == nil {
			firstErr = invokeErr
		}
		cancelParser()
	} else if fulfillErr := lifecycleSupervisor.WaitFulfilled(parserInvokeCtx, []string{parserRoleID}); fulfillErr != nil {
		forwarder.sendEvent("supervisor/lifecycle", decision.Event{Observer: "supervisor:" + container.Name, Event: "local_lifecycle_sigterm_fallback_used", Outcome: "broken", Peer: parserRoleID, Detail: fulfillErr.Error()})
		if firstErr == nil {
			firstErr = fulfillErr
		}
		cancelParser()
	}
	parserErr := <-parserErrs
	if firstErr == nil && parserErr != nil {
		firstErr = parserErr
	}
	kernelInvokeCtx, kernelInvokeCancel := context.WithTimeout(lifecycleCtx, shutdownGrace)
	defer kernelInvokeCancel()
	if invokeErr := lifecycleSupervisor.Invoke(kernelInvokeCtx, kernelRoleID, "parser_complete", ""); invokeErr != nil {
		forwarder.sendEvent("supervisor/lifecycle", decision.Event{Observer: "supervisor:" + container.Name, Event: "local_lifecycle_sigterm_fallback_used", Outcome: "broken", Peer: kernelRoleID, Detail: invokeErr.Error()})
		if firstErr == nil {
			firstErr = invokeErr
		}
		cancelKernel()
	} else if fulfillErr := lifecycleSupervisor.WaitFulfilled(kernelInvokeCtx, []string{kernelRoleID}); fulfillErr != nil {
		forwarder.sendEvent("supervisor/lifecycle", decision.Event{Observer: "supervisor:" + container.Name, Event: "local_lifecycle_sigterm_fallback_used", Outcome: "broken", Peer: kernelRoleID, Detail: fulfillErr.Error()})
		if firstErr == nil {
			firstErr = fulfillErr
		}
		cancelKernel()
	}
	kernelErr := <-kernelErrs
	if firstErr == nil && kernelErr != nil {
		firstErr = kernelErr
	}
	if forwardErr := forwarder.err(); firstErr == nil && forwardErr != nil {
		firstErr = forwardErr
	}
	return firstErr
}

func runKernelProcess(ctx context.Context, forwarder *processForwarder, shutdownGrace time.Duration, kernelBinary, configPath, containerName string, extraEnv []string, stdinSink func(io.WriteCloser)) error {
	// Intent: Start exactly one local kernel per container. It is transport
	// plumbing for app promises, not a trust authority or workflow owner.
	// Source: DI-galin
	return runProcess(ctx, forwarder, "kernel:"+containerName, shutdownGrace, kernelBinary, extraEnv, stdinSink, "-config", configPath, "-container", containerName)
}

func runParserRoleProcess(ctx context.Context, forwarder *processForwarder, shutdownGrace time.Duration, parserRoleBinary, configPath, containerName string, extraEnv []string, stdinSink func(io.WriteCloser)) error {
	// Intent: Start one local parser/builder role per container between apps and
	// the transport kernel so pCID-owned payload semantics stay out of the
	// transport kernel. Source: DI-gazin
	return runProcess(ctx, forwarder, "parser:"+containerName, shutdownGrace, parserRoleBinary, extraEnv, stdinSink, "-config", configPath, "-container", containerName)
}

func runAgentProcess(ctx context.Context, forwarder *processForwarder, shutdownGrace time.Duration, agentBinary, configPath, agentName string, extraEnv []string, stdinSink func(io.WriteCloser)) error {
	// Intent: A container supervisor starts independent local app processes
	// without sharing decision state; stdout/stderr pass through so the run log
	// remains auditable from Docker output. Source: DI-galin
	return runProcess(ctx, forwarder, "agent:"+agentName, shutdownGrace, agentBinary, extraEnv, stdinSink, "-config", configPath, "-node", agentName)
}

func runProcess(ctx context.Context, forwarder *processForwarder, roleName string, shutdownGrace time.Duration, binary string, extraEnv []string, stdinSink func(io.WriteCloser), args ...string) error {
	// Intent: Child processes still write ordinary stdout/stderr, but the
	// supervisor copies JSON event records to the observer-only collector instead
	// of relying on a Docker volume visible to agents. Source: DI-dirat
	command := exec.Command(binary, args...)
	command.Env = append(os.Environ(), extraEnv...)
	stdin, stdinErr := command.StdinPipe()
	if stdinErr != nil {
		return stdinErr
	}
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
	if stdinSink != nil {
		stdinSink(stdin)
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go forwarder.copyOutput(&waitGroup, roleName, "stdout", stdout, os.Stdout)
	go forwarder.copyOutput(&waitGroup, roleName, "stderr", stderr, os.Stderr)
	waitDone := make(chan error, 1)
	go func() {
		waitDone <- command.Wait()
	}()
	select {
	case waitErr := <-waitDone:
		waitGroup.Wait()
		closeProcessStdin(forwarder, roleName, stdin)
		if waitErr != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("%s failed: %w", roleName, waitErr)
		}
		return nil
	case <-ctx.Done():
		// Intent: Clean-run session accounting depends on kernels observing
		// shutdown and closing peer streams themselves. `exec.CommandContext`
		// kills children immediately, which hides terminal session records, so the
		// supervisor first sends SIGTERM and only kills after a bounded grace
		// window. Source: DI-homuj
		signalErr := command.Process.Signal(syscall.SIGTERM)
		if signalErr != nil && !errors.Is(signalErr, os.ErrProcessDone) {
			return fmt.Errorf("signal %s for shutdown: %w", roleName, signalErr)
		}
		timer := time.NewTimer(shutdownGrace)
		defer timer.Stop()
		select {
		case <-waitDone:
		case <-timer.C:
			killErr := command.Process.Kill()
			if killErr != nil && !errors.Is(killErr, os.ErrProcessDone) {
				return fmt.Errorf("kill %s after shutdown grace: %w", roleName, killErr)
			}
			<-waitDone
		}
		waitGroup.Wait()
		closeProcessStdin(forwarder, roleName, stdin)
		return nil
	}
}

func closeProcessStdin(forwarder *processForwarder, roleName string, stdin io.Closer) {
	// Intent: Victor's stdio lifecycle profile receives the supervisor's signed
	// shutdown token over stdin, then exits and may close its pipe before the
	// supervisor performs generic process cleanup. That is successful lifecycle
	// fulfillment, not a broken process promise. Source: DI-jafoj
	if closeErr := stdin.Close(); closeErr != nil && !isExpectedClosedStdin(closeErr) {
		forwarder.rememberErr(fmt.Errorf("close stdin for %s: %w", roleName, closeErr))
	}
}

func isExpectedClosedStdin(closeErr error) bool {
	if errors.Is(closeErr, os.ErrClosed) {
		return true
	}
	return strings.Contains(closeErr.Error(), "file already closed")
}

func lifecycleProfileForAgent(agent config.AgentConfig) string {
	// Intent: POC16 proves all three local_lifecycle_v1 transport profiles in one
	// clean run: ordinary roles use dedicated lifecycle TCP, parser roles are
	// invoked through their parser path, and Victor's stdio adapter receives
	// lifecycle invocation through stdin. Source: DI-jafoj
	if agent.Kind == "stdio_agent" {
		return protocol.LifecycleChannelStdio
	}
	return protocol.LifecycleChannelTCP
}

func parserRoleAddress(cfg config.Config, containerName string) (string, bool) {
	port, ok := cfg.ParserRoleAppPortForContainer(containerName)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("127.0.0.1:%d", port), true
}

func closeLifecycleSupervisor(supervisor *lifecycle.Supervisor) {
	if closeErr := supervisor.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
		fmt.Fprintf(os.Stderr, "close lifecycle supervisor: %v\n", closeErr)
	}
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
	var record eventstream.Record
	if unmarshalErr := json.Unmarshal([]byte(line), &record); unmarshalErr == nil && record.Kind == eventstream.KindMessageArtifact && record.MessageArtifact != nil {
		// Intent: Raw message artifacts are observer-only run review records.
		// The supervisor forwards them to the collector without interpreting the
		// envelope bytes or exposing the observer volume to app processes. Source:
		// DI-tuhop
		record.Source = forwarder.containerName + "/" + roleName + "/" + streamName
		sendErr := forwarder.client.Send(record)
		if sendErr != nil {
			forwarder.rememberErr(sendErr)
		}
		return
	}
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

func (forwarder *processForwarder) sendEvent(sourceSuffix string, event decision.Event) {
	sendErr := forwarder.client.Send(eventstream.Record{
		Kind:   eventstream.KindEvent,
		Source: forwarder.containerName + "/" + sourceSuffix,
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
		return "poc16-relationship-agent", nil
	case "fulfillment":
		return "poc16-fulfillment", nil
	case "postal_scale":
		return "poc16-postal-scale", nil
	case "ups_label_printer":
		return "poc16-ups-label-printer", nil
	case "printer_port":
		return "poc16-printer-port", nil
	case "accounting":
		return "poc16-accounting", nil
	case "wasm_agent":
		// Intent: POC16 adds heterogeneous runtime-adapter app roles as separate
		// supervised processes while preserving the same local-kernel routing
		// model used by POC12/POC13 agents. Source: DI-linof
		return "poc16-wasm-agent", nil
	case "stdio_agent":
		// Intent: The stdio adapter is the supervised app process for a worker
		// that sends and receives envelopes only over stdin/stdout. Source:
		// DI-linof
		return "poc16-stdio-adapter", nil
	default:
		return "", fmt.Errorf("agent %q has unsupported app kind %q", agentName, agent.Kind)
	}
}
