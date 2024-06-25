// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package gatherer

import (
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/mlc"

	"github.com/onmetal/inventory/pkg/block"
	"github.com/onmetal/inventory/pkg/cpu"
	"github.com/onmetal/inventory/pkg/distro"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/inventory"
	"github.com/onmetal/inventory/pkg/ipmi"
	"github.com/onmetal/inventory/pkg/lldp"
	"github.com/onmetal/inventory/pkg/mem"
	"github.com/onmetal/inventory/pkg/netlink"
	"github.com/onmetal/inventory/pkg/nic"
	"github.com/onmetal/inventory/pkg/numa"
	"github.com/onmetal/inventory/pkg/pci"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/virt"
)

type Svc struct {
	printer *printer.Svc

	dmiSvc     *dmi.Svc
	numaSvc    *numa.Svc
	blockSvc   *block.Svc
	pciSvc     *pci.Svc
	cpuInfoSvc *cpu.InfoSvc
	memInfoSvc *mem.InfoSvc
	mlcPerfSvc *mlc.PerfSvc
	lldpSvc    *lldp.Svc
	nicSvc     *nic.Svc
	ipmiSvc    *ipmi.Svc
	netlinkSvc *netlink.Svc
	virtSvc    *virt.Svc
	hostSvc    *host.Svc
	distroSvc  *distro.Svc
}

func NewSvc(printer *printer.Svc, opts ...Option) *Svc {
	svc := &Svc{
		printer: printer,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (s *Svc) Gather() *inventory.Inventory {
	setters := []func(inventory *inventory.Inventory) error{
		s.SetDMI,
		s.SetCPUInfo,
		s.SetMemInfo,
		// TODO Not gathering atm on regular run
		// mlc binary is not included as a dependency yet
		// Uncomment if dependency is met and benchmarking is required on regular run
		// s.SetMlcPerf,
		s.SetNumaNodes,
		s.SetBlockDevices,
		s.SetPCIBusDevices,
		s.SetIPMIDevices,
		s.SetNICs,
		s.SetLLDPFrames,
		s.SetNDPFrames,
		s.SetVirt,
		s.SetHost,
		s.SetDistro,
	}

	return s.GatherInOrder(setters)
}

func (s *Svc) GatherInOrder(setters []func(inventory *inventory.Inventory) error) *inventory.Inventory {
	inv := &inventory.Inventory{}

	for _, setter := range setters {
		err := setter(inv)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to set value"))
		}
	}

	return inv
}

func (s *Svc) SetDMI(inv *inventory.Inventory) error {
	data, err := s.dmiSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get dmi data")
	}
	inv.DMI = data
	return nil
}

func (s *Svc) SetCPUInfo(inv *inventory.Inventory) error {
	data, err := s.cpuInfoSvc.GetInfo()
	if err != nil {
		return errors.Wrap(err, "unable to get proc data")
	}
	inv.CPUInfo = data
	return nil
}

func (s *Svc) SetMemInfo(inv *inventory.Inventory) error {
	data, err := s.memInfoSvc.GetInfo()
	if err != nil {
		return errors.Wrap(err, "unable to get proc data")
	}
	inv.MemInfo = data
	return nil
}

func (s *Svc) SetMlcPerf(inv *inventory.Inventory) error {
	data, err := s.mlcPerfSvc.GetInfo()
	if err != nil {
		return errors.Wrap(err, "unable to get mlc data")
	}
	inv.MlcPerf = data
	return nil
}

func (s *Svc) SetNumaNodes(inv *inventory.Inventory) error {
	data, err := s.numaSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get numa data")
	}
	inv.NumaNodes = data
	return nil
}

func (s *Svc) SetBlockDevices(inv *inventory.Inventory) error {
	data, err := s.blockSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get block data")
	}
	inv.BlockDevices = data
	return nil
}

func (s *Svc) SetPCIBusDevices(inv *inventory.Inventory) error {
	data, err := s.pciSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get pci data")
	}
	inv.PCIBusDevices = data
	return nil
}

func (s *Svc) SetIPMIDevices(inv *inventory.Inventory) error {
	data, err := s.ipmiSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get ipmi data")
	}
	inv.IPMIDevices = data
	return nil
}

func (s *Svc) SetNICs(inv *inventory.Inventory) error {
	data, err := s.nicSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get nic data")
	}
	inv.NICs = data
	return nil
}

func (s *Svc) SetLLDPFrames(inv *inventory.Inventory) error {
	data, err := s.lldpSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get lldp data")
	}
	inv.LLDPFrames = data
	return nil
}

func (s *Svc) SetNDPFrames(inv *inventory.Inventory) error {
	data, err := s.netlinkSvc.GetIPv6NeighbourData()
	if err != nil {
		return errors.Wrap(err, "unable to get ndp data")
	}
	inv.NDPFrames = data
	return nil
}

func (s *Svc) SetVirt(inv *inventory.Inventory) error {
	data, err := s.virtSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get virt")
	}
	inv.Virtualization = data
	return nil
}

func (s *Svc) SetHost(inv *inventory.Inventory) error {
	hostInfo, err := s.hostSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get host info")
	}
	inv.Host = hostInfo
	return nil
}

func (s *Svc) SetDistro(inv *inventory.Inventory) error {
	if inv.Host == nil {
		cause := errors.New("no host data")
		return errors.Wrap(cause, "unable to get distro info")
	}

	distroInfo, err := s.distroSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get distro info")
	}
	inv.Distro = distroInfo
	return nil
}
