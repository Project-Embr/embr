package main

import (
	"io"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	guuid "github.com/google/uuid"
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

	fcBinary, err := exec.LookPath("firecracker")
	if err != nil {
		log.Errorf("No Firecracker Binary Found In PATH")
		return firecracker.Config{}, err
	}
	jailerBinary, err := exec.LookPath("jailer")
	if err != nil {
		log.Errorf("No Jailer Binary Found In PATH")
		return firecracker.Config{}, err
	}
	USER, err := user.Current()
	if err != nil {
		return firecracker.Config{}, err
	}
	chrootDir := "/machines"
	cgroupver := "2"
	Gid, err := strconv.Atoi(USER.Gid)
	Uid, err := strconv.Atoi(USER.Uid)
	jailer := &firecracker.JailerConfig{
		GID:            firecracker.Int(Gid),
		UID:            firecracker.Int(Uid),
		ID:             guuid.NewString(),
		ExecFile:       fcBinary,
		Daemonize:      true,
		NumaNode:       firecracker.Int(0),
		JailerBinary:   jailerBinary,
		ChrootBaseDir:  (chrootDir),
		CgroupVersion:  cgroupver,
		ChrootStrategy: firecracker.NewNaiveChrootStrategy(opts.FcKernelImage),
	}
	return firecracker.Config{
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
		}, JailerCfg: jailer,
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
