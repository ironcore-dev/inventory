package inventory

import (
	"github.com/onmetal/inventory/pkg/block"
	"github.com/onmetal/inventory/pkg/cpu"
	"github.com/onmetal/inventory/pkg/distro"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/ipmi"
	"github.com/onmetal/inventory/pkg/lldp/frame"
	"github.com/onmetal/inventory/pkg/mem"
	"github.com/onmetal/inventory/pkg/netlink"
	"github.com/onmetal/inventory/pkg/nic"
	"github.com/onmetal/inventory/pkg/numa"
	"github.com/onmetal/inventory/pkg/pci"
	"github.com/onmetal/inventory/pkg/virt"
)

type Inventory struct {
	DMI            *dmi.DMI
	MemInfo        *mem.Info
	CPUInfo        []cpu.Info
	NumaNodes      []numa.Node
	BlockDevices   []block.Device
	PCIBusDevices  []pci.Bus
	IPMIDevices    []ipmi.Device
	NICs           []nic.Device
	LLDPFrames     []frame.Frame
	NDPFrames      []netlink.IPv6Neighbour
	Virtualization *virt.Virtualization
	Host           *host.Info
	Distro         *distro.Distro
}
