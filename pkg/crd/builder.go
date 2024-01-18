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

package crd

import (
	"sort"
	"strconv"
	"strings"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/onmetal/inventory/pkg/inventory"
	"github.com/onmetal/inventory/pkg/lldp/frame"
	"github.com/onmetal/inventory/pkg/netlink"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/utils"
)

const (
	CLoopbackNICPrefix = "lo"
	CDockerNICPrefix   = "docker"
)

var CNICPrefixesToExclude = []string{
	CLoopbackNICPrefix,
	CDockerNICPrefix,
}

type BuilderSvc struct {
	printer *printer.Svc
}

func NewBuilderSvc(printer *printer.Svc) *BuilderSvc {
	return &BuilderSvc{
		printer: printer,
	}
}

func (s *BuilderSvc) Build(inv *inventory.Inventory) (*metalv1alpha4.Inventory, error) {
	setters := []func(*metalv1alpha4.Inventory, *inventory.Inventory){
		s.SetSystem,
		s.SetIPMIs,
		s.SetBlocks,
		s.SetMemory,
		s.SetCPUs,
		s.SetNUMANodes,
		s.SetPCIDevices,
		s.SetNICs,
		s.SetVirt,
		s.SetHost,
		s.SetDistro,
	}

	return s.BuildInOrder(inv, setters)
}

func (s *BuilderSvc) BuildInOrder(inv *inventory.Inventory, setters []func(*metalv1alpha4.Inventory, *inventory.Inventory)) (*metalv1alpha4.Inventory, error) {
	cr := &metalv1alpha4.Inventory{
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       metalv1alpha4.InventorySpec{},
	}

	for _, setter := range setters {
		setter(cr, inv)
	}

	return cr, nil
}

func (s *BuilderSvc) SetSystem(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if inv.DMI == nil {
		return
	}

	dmi := inv.DMI
	if dmi.SystemInformation == nil {
		return
	}

	if inv.Host == nil {
		return
	}

	// SONiC switches has dumb UUIDs like 03000200-0400-0500-0006-000700080009, maybe
	// the same on any switch, so it was decided to use md5 hash of serial number as UUID
	hostUUID := dmi.SystemInformation.UUID
	if inv.Host.Type == utils.CSwitchType {
		hostUUID = getUUID(CSonicNamespace, dmi.SystemInformation.SerialNumber)
	}
	cr.Name = hostUUID

	cr.Spec.System = &metalv1alpha4.SystemSpec{
		ID:           hostUUID,
		Manufacturer: dmi.SystemInformation.Manufacturer,
		ProductSKU:   dmi.SystemInformation.SKUNumber,
		SerialNumber: dmi.SystemInformation.SerialNumber,
	}
}

func (s *BuilderSvc) SetIPMIs(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	ipmiDevCount := len(inv.IPMIDevices)
	if ipmiDevCount == 0 {
		return
	}

	ipmis := make([]metalv1alpha4.IPMISpec, ipmiDevCount)

	for i, ipmiDev := range inv.IPMIDevices {
		ipmi := metalv1alpha4.IPMISpec{
			IPAddress:  ipmiDev.IPAddress,
			MACAddress: ipmiDev.MACAddress,
		}

		ipmis[i] = ipmi
	}

	sort.Slice(ipmis, func(i, j int) bool {
		iStr := ipmis[i].MACAddress + ipmis[i].IPAddress
		jStr := ipmis[j].MACAddress + ipmis[j].IPAddress
		return iStr < jStr
	})

	cr.Spec.IPMIs = ipmis
}

func (s *BuilderSvc) SetBlocks(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if len(inv.BlockDevices) == 0 {
		return
	}

	blocks := make([]metalv1alpha4.BlockSpec, 0)
	var capacity uint64 = 0

	for _, blockDev := range inv.BlockDevices {
		// Filter non physical devices
		if blockDev.Type == "" {
			continue
		}

		var partitionTable *metalv1alpha4.PartitionTableSpec
		if blockDev.PartitionTable != nil {
			table := blockDev.PartitionTable

			partitionTable = &metalv1alpha4.PartitionTableSpec{
				Type: string(table.Type),
			}

			partCount := len(table.Partitions)
			if partCount > 0 {
				parts := make([]metalv1alpha4.PartitionSpec, partCount)

				for i, partition := range table.Partitions {
					part := metalv1alpha4.PartitionSpec{
						ID:   partition.ID,
						Name: partition.Name,
						Size: partition.Size,
					}

					parts[i] = part
				}

				sort.Slice(parts, func(i, j int) bool {
					return parts[i].ID < parts[j].ID
				})

				partitionTable.Partitions = parts
			}
		}

		block := metalv1alpha4.BlockSpec{
			Name:           blockDev.Name,
			Type:           blockDev.Type,
			Rotational:     blockDev.Rotational,
			Bus:            blockDev.Vendor,
			Model:          blockDev.Model,
			Size:           blockDev.Size,
			PartitionTable: partitionTable,
		}

		capacity += blockDev.Size
		blocks = append(blocks, block)
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Name < blocks[j].Name
	})

	cr.Spec.Blocks = blocks
}

func (s *BuilderSvc) SetMemory(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if inv.MemInfo == nil {
		return
	}

	cr.Spec.Memory = &metalv1alpha4.MemorySpec{
		Total: inv.MemInfo.MemTotal,
	}
}

func (s *BuilderSvc) SetMLCPerf(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	// TODO set data to inventory when CRD will get perf fields
}

func (s *BuilderSvc) SetCPUs(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if len(inv.CPUInfo) == 0 {
		return
	}

	cpuMarkMap := make(map[uint64]metalv1alpha4.CPUSpec)

	for _, cpuInfo := range inv.CPUInfo {
		if val, ok := cpuMarkMap[cpuInfo.PhysicalID]; ok {
			val.LogicalIDs = append(val.LogicalIDs, cpuInfo.Processor)
			cpuMarkMap[cpuInfo.PhysicalID] = val
			continue
		}

		cpu := metalv1alpha4.CPUSpec{
			PhysicalID:      cpuInfo.PhysicalID,
			LogicalIDs:      []uint64{cpuInfo.Processor},
			Cores:           cpuInfo.CpuCores,
			Siblings:        cpuInfo.Siblings,
			VendorID:        cpuInfo.VendorID,
			Family:          cpuInfo.CPUFamily,
			Model:           cpuInfo.Model,
			ModelName:       cpuInfo.ModelName,
			Stepping:        cpuInfo.Stepping,
			Microcode:       cpuInfo.Microcode,
			MHz:             resource.MustParse(cpuInfo.CPUMHz),
			CacheSize:       cpuInfo.CacheSize,
			FPU:             cpuInfo.FPU,
			FPUException:    cpuInfo.FPUException,
			CPUIDLevel:      cpuInfo.CPUIDLevel,
			WP:              cpuInfo.WP,
			Flags:           cpuInfo.Flags,
			VMXFlags:        cpuInfo.VMXFlags,
			Bugs:            cpuInfo.Bugs,
			BogoMIPS:        resource.MustParse(cpuInfo.BogoMIPS),
			CLFlushSize:     cpuInfo.CLFlushSize,
			CacheAlignment:  cpuInfo.CacheAlignment,
			AddressSizes:    cpuInfo.AddressSizes,
			PowerManagement: cpuInfo.PowerManagement,
		}
		sort.Strings(cpu.Flags)
		sort.Strings(cpu.VMXFlags)
		sort.Strings(cpu.Bugs)
		sort.Slice(cpu.LogicalIDs, func(i, j int) bool {
			return cpu.LogicalIDs[i] < cpu.LogicalIDs[j]
		})

		cpuMarkMap[cpuInfo.PhysicalID] = cpu
	}

	cpus := make([]metalv1alpha4.CPUSpec, 0)
	for _, v := range cpuMarkMap {
		cpus = append(cpus, v)
	}

	sort.Slice(cpus, func(i, j int) bool {
		return cpus[i].PhysicalID < cpus[j].PhysicalID
	})

	cr.Spec.CPUs = cpus
}

func (s *BuilderSvc) SetNUMANodes(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if len(inv.NumaNodes) == 0 {
		return
	}

	numaNodes := make([]metalv1alpha4.NumaSpec, len(inv.NumaNodes))
	for idx, numaNode := range inv.NumaNodes {
		numaNodes[idx] = metalv1alpha4.NumaSpec{
			ID:        numaNode.ID,
			CPUs:      numaNode.CPUs,
			Distances: numaNode.Distances,
		}
		if numaNode.Memory != nil {
			numaNodes[idx].Memory = &metalv1alpha4.MemorySpec{
				Total: numaNode.Memory.MemTotal,
			}
		}
	}

	sort.Slice(numaNodes, func(i, j int) bool {
		return numaNodes[i].ID < numaNodes[j].ID
	})

	cr.Spec.NUMA = numaNodes
}

func (s *BuilderSvc) SetPCIDevices(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if len(inv.PCIBusDevices) == 0 {
		return
	}

	pciDevices := make([]metalv1alpha4.PCIDeviceSpec, 0)
	for _, pciBus := range inv.PCIBusDevices {
		for _, pciDevice := range pciBus.Devices {
			pciDeviceSpec := metalv1alpha4.PCIDeviceSpec{
				BusID:   pciBus.ID,
				Address: pciDevice.Address,
			}
			if pciDevice.Vendor != nil {
				pciDeviceSpec.Vendor = &metalv1alpha4.PCIDeviceDescriptionSpec{
					ID:   pciDevice.Vendor.ID,
					Name: pciDevice.Vendor.Name,
				}
			}
			if pciDevice.Subvendor != nil {
				pciDeviceSpec.Subvendor = &metalv1alpha4.PCIDeviceDescriptionSpec{
					ID:   pciDevice.Subvendor.ID,
					Name: pciDevice.Subvendor.Name,
				}
			}
			if pciDevice.Type != nil {
				pciDeviceSpec.Type = &metalv1alpha4.PCIDeviceDescriptionSpec{
					ID:   pciDevice.Type.ID,
					Name: pciDevice.Type.Name,
				}
			}
			if pciDevice.Subtype != nil {
				pciDeviceSpec.Subtype = &metalv1alpha4.PCIDeviceDescriptionSpec{
					ID:   pciDevice.Subtype.ID,
					Name: pciDevice.Subtype.Name,
				}
			}
			if pciDevice.Class != nil {
				pciDeviceSpec.Class = &metalv1alpha4.PCIDeviceDescriptionSpec{
					ID:   pciDevice.Class.ID,
					Name: pciDevice.Class.Name,
				}
			}
			if pciDevice.Subclass != nil {
				pciDeviceSpec.Subclass = &metalv1alpha4.PCIDeviceDescriptionSpec{
					ID:   pciDevice.Subclass.ID,
					Name: pciDevice.Subclass.Name,
				}
			}
			if pciDevice.ProgrammingInterface != nil {
				pciDeviceSpec.ProgrammingInterface = &metalv1alpha4.PCIDeviceDescriptionSpec{
					ID:   pciDevice.ProgrammingInterface.ID,
					Name: pciDevice.ProgrammingInterface.Name,
				}
			}

			pciDevices = append(pciDevices, pciDeviceSpec)
		}
	}

	sort.Slice(pciDevices, func(i, j int) bool {
		return pciDevices[i].Address < pciDevices[j].Address
	})

	cr.Spec.PCIDevices = pciDevices
}

func (s *BuilderSvc) SetNICs(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if len(inv.NICs) == 0 {
		return
	}

	if inv.Host == nil {
		return
	}

	lldpMap := make(map[int][]metalv1alpha4.LLDPSpec)
	for _, f := range inv.LLDPFrames {
		checkMap := make(map[frame.Capability]struct{})
		enabledCapabilities := make([]metalv1alpha4.LLDPCapabilities, 0)
		for _, capability := range f.EnabledCapabilities {
			if _, ok := checkMap[capability]; !ok {
				enabledCapabilities = append(enabledCapabilities, metalv1alpha4.LLDPCapabilities(capability))
				checkMap[capability] = struct{}{}
			}
		}
		sort.Slice(enabledCapabilities, func(i, j int) bool {
			return enabledCapabilities[i] < enabledCapabilities[j]
		})
		id, _ := strconv.Atoi(f.InterfaceID)
		l := metalv1alpha4.LLDPSpec{
			ChassisID:         f.ChassisID,
			SystemName:        f.SystemName,
			SystemDescription: f.SystemDescription,
			PortID:            f.PortID,
			PortDescription:   f.PortDescription,
			Capabilities:      enabledCapabilities,
		}

		if _, ok := lldpMap[id]; !ok {
			lldpMap[id] = make([]metalv1alpha4.LLDPSpec, 0)
		}

		lldpMap[id] = append(lldpMap[id], l)
	}

	ndpMap := make(map[int][]metalv1alpha4.NDPSpec)
	for _, ndp := range inv.NDPFrames {
		// filtering no arp as ip neigh does
		if ndp.State == netlink.CNeighbourNoARPCacheState {
			continue
		}

		n := metalv1alpha4.NDPSpec{
			IPAddress:  ndp.IP,
			MACAddress: ndp.MACAddress,
			State:      string(ndp.State),
		}

		if _, ok := ndpMap[ndp.DeviceIndex]; !ok {
			ndpMap[ndp.DeviceIndex] = make([]metalv1alpha4.NDPSpec, 0)
		}

		ndpMap[ndp.DeviceIndex] = append(ndpMap[ndp.DeviceIndex], n)
	}

	nics := make([]metalv1alpha4.NICSpec, 0)
	for _, nic := range inv.NICs {
		// it was reported that loopback and docker interfaces on some systems may have
		// PCI address assigned (but they have weird format), so we need to exclude them too
		shouldExclude := false
		for _, prefix := range CNICPrefixesToExclude {
			if strings.HasPrefix(nic.Name, prefix) {
				shouldExclude = true
				break
			}
		}

		if shouldExclude {
			continue
		}

		// filter non-physical interfaces according to type of inventorying host
		if inv.Host.Type == utils.CSwitchType {
			if nic.PCIAddress == "" && !strings.HasPrefix(nic.Name, "Ethernet") {
				continue
			}
		} else {
			if nic.PCIAddress == "" {
				continue
			}
		}

		lldps := lldpMap[int(nic.InterfaceIndex)]
		sort.Slice(lldps, func(i, j int) bool {
			iStr := lldps[i].ChassisID + lldps[i].PortID
			jStr := lldps[j].ChassisID + lldps[j].PortID
			return iStr < jStr
		})
		ndps := ndpMap[int(nic.InterfaceIndex)]
		sort.Slice(ndps, func(i, j int) bool {
			iStr := ndps[i].MACAddress + ndps[i].IPAddress
			jStr := ndps[j].MACAddress + ndps[j].IPAddress
			return iStr < jStr
		})

		ns := metalv1alpha4.NICSpec{
			Name:       nic.Name,
			PCIAddress: nic.PCIAddress,
			MACAddress: nic.Address,
			MTU:        nic.MTU,
			Speed:      nic.Speed,
			LLDPs:      lldps,
			NDPs:       ndps,
			ActiveFEC:  nic.FEC,
			Lanes:      nic.Lanes,
		}

		nics = append(nics, ns)
	}

	sort.Slice(nics, func(i, j int) bool {
		return nics[i].Name < nics[j].Name
	})

	cr.Spec.NICs = nics
}

func (s *BuilderSvc) SetVirt(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if inv.Virtualization == nil {
		return
	}

	cr.Spec.Virt = &metalv1alpha4.VirtSpec{
		VMType: string(inv.Virtualization.Type),
	}
}

func (s *BuilderSvc) SetHost(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if inv.Host == nil {
		return
	}
	cr.Spec.Host = &metalv1alpha4.HostSpec{
		Name: inv.Host.Name,
	}
}

func (s *BuilderSvc) SetDistro(cr *metalv1alpha4.Inventory, inv *inventory.Inventory) {
	if inv.Distro == nil {
		return
	}
	cr.Spec.Distro = &metalv1alpha4.DistroSpec{
		BuildVersion:  inv.Distro.BuildVersion,
		DebianVersion: inv.Distro.DebianVersion,
		KernelVersion: inv.Distro.KernelVersion,
		AsicType:      inv.Distro.AsicType,
		CommitID:      inv.Distro.CommitID,
		BuildDate:     inv.Distro.BuildDate,
		BuildNumber:   inv.Distro.BuildNumber,
		BuildBy:       inv.Distro.BuildBy,
	}
}
