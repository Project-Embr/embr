package main

import (
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"os/exec"
	"os/user"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

func newOptions() *options {
	return &options{}
}

type options struct {
	FcKernelImage   string   `description:"Kernel image path"`
	FcRootDrivePath string   `description:"RootFS path"`
	VCpuCount       int64    `json:"VCpuCount,omitempty"`
	MemSizeMib      int64    `json:"MemSizeMib,omitempty"`
	CNIConfigPath   string   `discription:"CNI network configuration path"`
	CNIPluginsPath  []string `discription:"CNI plugins path"`
	CNINetnsPath    string   `discription:"CNI Netns path"`

	Id           string `long:"id" description:"Jailer VMM id"`

	Uid      int `long:"uid" description:"Jailer uid"`
	Gid      int `long:"gid" description:"Jailer gid"`
	NumaNode int `long:"node" description:"Jailer numa node"`
}

// Converts options to a usable firecracker config
func (opts *options) createFirecrackerConfig() (firecracker.Config, error, string) {
	// setup NICs
	var NICs []firecracker.NetworkInterface
	// BlockDevices
	blockDevices, err := opts.getBlockDevices()
	if err != nil {
		return firecracker.Config{}, err, ""
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

	// Append socket path to Netns
	opts.CNINetnsPath = strings.Join([]string{
		opts.CNINetnsPath,
		socketPath,
	},
		"",
	)


	
	//err will need to be handled
	fcBinary, err := exec.LookPath("firecracker")
	if err != nil{
		return firecracker.Config{}, err, ""
	}
	jailerBinary, err := exec.LookPath("jailer")
	if err != nil{
		return firecracker.Config{}, err, ""
	}
	USER, err := user.Current() 
	if err != nil{
		return firecracker.Config{}, err, ""
	}
	Gid, err := strconv.Atoi(USER.Gid)
	Uid, err := strconv.Atoi(USER.Uid)
	//Creates the jailer
	jailer := &firecracker.JailerConfig{
		GID:            firecracker.Int(Gid),
		UID:            firecracker.Int(Uid),
		ID:             opts.Id,
		NumaNode:       firecracker.Int(opts.NumaNode),
		ExecFile:       fcBinary,
		JailerBinary:   jailerBinary,
		Daemonize:     	true,
		ChrootStrategy: firecracker.NewNaiveChrootStrategy(opts.FcKernelImage),
		Stdout:         os.Stdout,
		Stderr:         os.Stderr,
		Stdin:          os.Stdin,
	}
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
			VcpuCount:  firecracker.Int64(opts.VCpuCount),
			Smt:        firecracker.Bool(true),
			MemSizeMib: firecracker.Int64(opts.MemSizeMib),
		},
		JailerCfg: jailer,
	}, nil, socketPath
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
