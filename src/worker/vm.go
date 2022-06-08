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

//Creates a new firecracker VM and returns the command channel
func createNewVM(opts *options) chan string{
	command := make(chan string, 1)
	err := make(chan error, 1)
	go runVM(context.Background(), opts, err, command)

	if (<- err == nil){
		fmt.Println("machine started successfully")
	} else{
		fmt.Println("Failed to create machine")
	}
	return command
}

// Run a firecracker VM
func runVM(ctx context.Context, opts *options, er chan<- error, cmd <-chan string) error {
	// options -> firecracker config
	fcCfg, err, socketPath := opts.createFirecrackerConfig()
	//if err != nil {
	//	log.Errorf("Error: %s", err)
	//	return err
	//}
	//logger := log.New()

	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()
	//machineOpts := []firecracker.Opt{
	//	firecracker.WithLogger(log.NewEntry(logger)),
	//}

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

	if err := machine.Start(vmmCtx); err != nil {
		er <- err
		return fmt.Errorf("failed to start machine: %v", err)
	}

	signalHandlers(vmmCtx, machine)
	er <- nil
	// Deleted signal handlers, re add
		// Wait for a signal and triger stopVMM
	if(<- cmd == "shutdown"){
		machine.Shutdown(ctx)
	}
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