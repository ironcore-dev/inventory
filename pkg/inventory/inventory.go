// Copyright 2023 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	"github.com/onmetal/inventory/pkg/mlc"
	"github.com/onmetal/inventory/pkg/netlink"
	"github.com/onmetal/inventory/pkg/nic"
	"github.com/onmetal/inventory/pkg/numa"
	"github.com/onmetal/inventory/pkg/pci"
	"github.com/onmetal/inventory/pkg/virt"
)

type Inventory struct {
	DMI            *dmi.DMI
	MemInfo        *mem.Info
	MlcPerf        *mlc.Perf
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
