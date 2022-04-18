package main

import (
	"github.com/firecracker-microvm/firecracker-go-sdk"
)

func cni(machine *firecracker.Machine) {
	machine.Cfg.NetworkInterfaces = []firecracker.NetworkInterface{{
		CNIConfiguration: &firecracker.CNIConfiguration{
			NetworkName: "fcnet",
			IfName:      "veth0",
			ConfDir:     "cni/conf.d",
			BinPath:     []string{"../../plugins/bin"},
		},
	}}
	machine.Cfg.NetNS = "ext/netns"
}
