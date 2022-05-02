package main

import (
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

func newOptions() *options {
	return &options{}
}

type options struct {
	FcKernelImage   string   `description:"Kernel image path"`
	FcRootDrivePath string   `description:"RootFS path"`
	CNIConfigPath   string   `discription:"CNI network configuration path"`
	CNIPluginsPath  []string `discription:"CNI plugins path"`
}

// Converts options to a usable firecracker config
func (opts *options) createFirecrackerConfig() (firecracker.Config, error) {
	// setup NICs
	var NICs []firecracker.NetworkInterface

	// BlockDevices
	blockDevices, err := opts.getBlockDevices()
	if err != nil {
		return firecracker.Config{}, err
	}

	// vsocks
	var vsocks []firecracker.VsockDevice

	// fifos
	var fifo io.WriteCloser

	// Generate socket path
	var socketPath = strings.Join([]string{
		".firecracker.sock",
		strconv.Itoa(os.Getpid()),
		strconv.Itoa(rand.Intn(1000))},
		"-",
	)

	return firecracker.Config{
		SocketPath:        socketPath,
		LogLevel:          "Debug",
		FifoLogWriter:     fifo,
		KernelImagePath:   opts.FcKernelImage,
		KernelArgs:        "ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules",
		Drives:            blockDevices,
		NetworkInterfaces: NICs,
		VsockDevices:      vsocks,
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(1),
			Smt:        firecracker.Bool(true),
			MemSizeMib: firecracker.Int64(512),
		},
	}, nil
}

// constructs a list of drives from the options config
func (opts *options) getBlockDevices() ([]models.Drive, error) {
	blockDevices := []models.Drive{}

	rootDrivePath, readOnly := parseDevice(opts.FcRootDrivePath)
	rootDrive := models.Drive{
		DriveID:      firecracker.String("1"),
		PathOnHost:   firecracker.String(rootDrivePath),
		IsReadOnly:   firecracker.Bool(readOnly),
		IsRootDevice: firecracker.Bool(true),
		Partuuid:     "",
	}
	blockDevices = append(blockDevices, rootDrive)
	return blockDevices, nil
}

// Check string for readonly marker
func parseDevice(entry string) (path string, readOnly bool) {
	if strings.HasSuffix(entry, ":ro") {
		return strings.TrimSuffix(entry, ":ro"), true
	}

	return strings.TrimSuffix(entry, ":rw"), false
}
