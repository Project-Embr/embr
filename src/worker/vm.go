package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	log "github.com/sirupsen/logrus"
)

// Run a firecracker VM
func runVM(ctx context.Context, opts *options) error {
	// options -> firecracker config
	fcCfg, err := opts.createFirecrackerConfig()
	if err != nil {
		log.Errorf("Error: %s", err)
		return err
	}
	logger := log.New()

	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()

	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(logger)),
	}

	var firecrackerBinary string
	firecrackerBinary, err = exec.LookPath("firecracker")
	if os.IsNotExist(err) {
		return fmt.Errorf("binary %q does not exist: %v", firecrackerBinary, err)
	}
	if err != nil {
		return fmt.Errorf("failed to start binary, %q: %v", firecrackerBinary, err)
	}

	machine, err := firecracker.NewMachine(vmmCtx, fcCfg, machineOpts...)
	if err != nil {
		return fmt.Errorf("failed creating machine: %s", err)
	}

	startVeth(machine, opts)

	if err := machine.Start(vmmCtx); err != nil {
		return fmt.Errorf("failed to start machine: %v", err)
	}

	// Stop VM when this function exits
	defer machine.StopVMM()

	signalHandlers(vmmCtx, machine)

	// wait for the VM to exit
	if err := machine.Wait(vmmCtx); err != nil {
		return fmt.Errorf("wait returned an error %s", err)
	}

	log.Printf("Machine started successfully")
	return nil
}

// Create custom signal handlers
func signalHandlers(ctx context.Context, m *firecracker.Machine) {
	go func() {
		// Reset signal handlers
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		channel := make(chan os.Signal, 1)
		signal.Notify(channel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		var signal os.Signal = <-channel
		if signal == syscall.SIGTERM || signal == os.Interrupt {
			log.Printf("Caught signal: %s, clean shutdown", signal.String())
			m.Shutdown(ctx)
		} else if signal == syscall.SIGQUIT {
			log.Printf("Caught signal: %s, force shutdown", signal.String())
			m.StopVMM()
		}
	}()
}
