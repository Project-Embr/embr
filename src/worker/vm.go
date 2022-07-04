package main

import (
	"context"
	"fmt"
	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	client "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Run a firecracker VM
func runVM(ctx context.Context, opts *options, er chan<- error, cmd chan string) error {

	fcCfg, err, socketPath := opts.createFirecrackerConfig()
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()

	var firecrackerBinary string
	firecrackerBinary, err = exec.LookPath("firecracker")
	if os.IsNotExist(err) {
		er <- err
		return fmt.Errorf("binary %q does not exist: %v", firecrackerBinary, err)
	}
	if err != nil {
		er <- err
		return fmt.Errorf("failed to start binary, %q: %v", firecrackerBinary, err)
	}

	c := firecracker.VMCommandBuilder{}.
		WithSocketPath(socketPath).
		WithBin(firecrackerBinary).
		Build(vmmCtx)

	machine, err := firecracker.NewMachine(vmmCtx, fcCfg, firecracker.WithProcessRunner(c))
	if err != nil {
		er <- err
		return fmt.Errorf("failed creating machine: %s", err)
	}

	startVeth(machine, opts)

	if err := machine.Start(vmmCtx); err != nil {
		er <- err
		return fmt.Errorf("failed to start machine: %v", err)
	}

	er <- nil

	if <-cmd == "shutdown" {
		machine.Shutdown(ctx)
		cmd <- "Shutdown Complete"
	}
	return nil
}

// Custom Signal Handlers
func SignalHandlers(etcdClient *client.Client, etcdServer *embed.Etcd, VMPointer *[]chan string) {
	go func() {
		// Reset signal handlers
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		channel := make(chan os.Signal, 1)
		signal.Notify(channel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		var signal os.Signal = <-channel
		VM := *VMPointer
		if signal == syscall.SIGTERM || signal == os.Interrupt {
			for i := 0; i < len(VM); i++ {
				VM[i] <- "shutdown"
				<-VM[i]
			}
			etcdClient.Close()
			etcdServer.Close()
		} else if signal == syscall.SIGQUIT {
			for i := 0; i < len(VM); i++ {
				VM[i] <- "shutdown"
				<-VM[i]
			}
			etcdClient.Close()
			etcdServer.Close()
		}
	}()
}
