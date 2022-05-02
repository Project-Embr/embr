package main

import "github.com/firecracker-microvm/firecracker-go-sdk"

func cni(machine *firecracker.Machine, opts *options) {
	machine.Cfg.NetworkInterfaces = []firecracker.NetworkInterface{{
		CNIConfiguration: &firecracker.CNIConfiguration{
			NetworkName: "fcnet",
			IfName:      "veth0",
			ConfDir:     opts.CNIConfigPath,
			BinPath:     opts.CNIPluginsPath,
		},
	}}
	machine.Cfg.NetNS = opts.CNINetnsPath
}
