package gatherer

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/block"
	"github.com/onmetal/inventory/pkg/cpu"
	"github.com/onmetal/inventory/pkg/crd"
	"github.com/onmetal/inventory/pkg/distro"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/flags"
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

const (
	COKRetCode  = 0
	CErrRetCode = -1
)

type Svc struct {
	printer *printer.Svc

	crdSvc *crd.Svc

	dmiSvc     *dmi.Svc
	numaSvc    *numa.Svc
	blockSvc   *block.Svc
	pciSvc     *pci.Svc
	cpuInfoSvc *cpu.InfoSvc
	memInfoSvc *mem.InfoSvc
	lldpSvc    *lldp.Svc
	nicSvc     *nic.Svc
	ipmiSvc    *ipmi.Svc
	netlinkSvc *netlink.Svc
	virtSvc    *virt.Svc
	hostSvc    *host.Svc
	distroSvc  *distro.Svc
}

func NewSvc() (*Svc, int) {
	f := flags.NewFlags()

	p := printer.NewSvc(f.Verbose)

	crdSvc, err := crd.NewSvc(f.Kubeconfig, f.KubeNamespace)
	if err != nil {
		p.Err(errors.Wrapf(err, "unable to create k8s resorce svc"))
		return nil, CErrRetCode
	}

	pciIDs, err := pci.NewIDs()
	if err != nil {
		p.Err(errors.Wrapf(err, "unable to load PCI IDs"))
		return nil, CErrRetCode
	}

	rawDmiSvc := dmi.NewRawSvc(f.Root)
	dmiSvc := dmi.NewSvc(p, rawDmiSvc)

	cpuInfoSvc := cpu.NewInfoSvc(p, f.Root)
	memInfoSvc := mem.NewInfoSvc(p, f.Root)

	numaStatSvc := numa.NewStatSvc(p)
	numaNodeSvc := numa.NewNodeSvc(memInfoSvc, numaStatSvc)
	numaSvc := numa.NewSvc(p, numaNodeSvc, f.Root)

	partitionTableSvc := block.NewPartitionTableSvc(f.Root)
	blockDeviceStatSvc := block.NewDeviceStatSvc(p)
	blockDeviceSvc := block.NewDeviceSvc(p, partitionTableSvc, blockDeviceStatSvc)
	blockSvc := block.NewSvc(p, blockDeviceSvc, f.Root)

	pciDevSvc := pci.NewDeviceSvc(p, pciIDs)
	pciBusSvc := pci.NewBusSvc(p, pciDevSvc)
	pciSvc := pci.NewSvc(p, pciBusSvc, f.Root)

	hostSvc := host.NewSvc(p, f.Root)

	redisSvc := lldp.NewRedisSvc(f.Root)
	lldpFrameInfoSvc := lldp.NewFrameSvc(p)
	lldpSvc := lldp.NewSvc(p, lldpFrameInfoSvc, hostSvc, redisSvc, f.Root)

	nicDevSvc := nic.NewDeviceSvc(p)
	nicSvc := nic.NewSvc(p, nicDevSvc, f.Root)

	ipmiDevInfoSvc := ipmi.NewDeviceSvc(p)
	ipmiSvc := ipmi.NewSvc(p, ipmiDevInfoSvc, f.Root)

	nlSvc := netlink.NewSvc(p, f.Root)

	virtSvc := virt.NewSvc(dmiSvc, cpuInfoSvc, f.Root)

	distroSvc := distro.NewSvc(p, hostSvc, f.Root)

	return &Svc{
		printer:    p,
		crdSvc:     crdSvc,
		dmiSvc:     dmiSvc,
		numaSvc:    numaSvc,
		blockSvc:   blockSvc,
		pciSvc:     pciSvc,
		cpuInfoSvc: cpuInfoSvc,
		memInfoSvc: memInfoSvc,
		lldpSvc:    lldpSvc,
		nicSvc:     nicSvc,
		ipmiSvc:    ipmiSvc,
		netlinkSvc: nlSvc,
		virtSvc:    virtSvc,
		hostSvc:    hostSvc,
		distroSvc:  distroSvc,
	}, 0
}

func (s *Svc) Gather() int {
	inv := &inventory.Inventory{}

	setters := []func(inventory *inventory.Inventory) error{
		s.setDMI,
		s.setCPUInfo,
		s.setMemInfo,
		s.setNumaNodes,
		s.setBlockDevices,
		s.setPCIBusDevices,
		s.setIPMIDevices,
		s.setNICs,
		s.setLLDPFrames,
		s.setNDPFrames,
		s.setVirt,
		s.setHost,
		s.setDistro,
	}

	for _, setter := range setters {
		err := setter(inv)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to set value"))
		}
	}

	jsonBytes, err := json.Marshal(inv)
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to marshal result to json"))
	}

	var prettifiedJsonBuf bytes.Buffer
	if err := json.Indent(&prettifiedJsonBuf, jsonBytes, "", "\t"); err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to indent json"))
	}

	s.printer.VOut("Gathered data:")
	s.printer.VOut(prettifiedJsonBuf.String())

	if err := s.crdSvc.BuildAndSave(inv); err != nil {
		s.printer.Err(errors.Wrap(err, "unable to save inventory resource"))
		return CErrRetCode
	}

	return COKRetCode
}

func (s *Svc) setDMI(inv *inventory.Inventory) error {
	data, err := s.dmiSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get dmi data")
	}
	inv.DMI = data
	return nil
}

func (s *Svc) setCPUInfo(inv *inventory.Inventory) error {
	data, err := s.cpuInfoSvc.GetInfo()
	if err != nil {
		return errors.Wrap(err, "unable to get proc data")
	}
	inv.CPUInfo = data
	return nil
}

func (s *Svc) setMemInfo(inv *inventory.Inventory) error {
	data, err := s.memInfoSvc.GetInfo()
	if err != nil {
		return errors.Wrap(err, "unable to get proc data")
	}
	inv.MemInfo = data
	return nil
}

func (s *Svc) setNumaNodes(inv *inventory.Inventory) error {
	data, err := s.numaSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get numa data")
	}
	inv.NumaNodes = data
	return nil
}

func (s *Svc) setBlockDevices(inv *inventory.Inventory) error {
	data, err := s.blockSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get block data")
	}
	inv.BlockDevices = data
	return nil
}

func (s *Svc) setPCIBusDevices(inv *inventory.Inventory) error {
	data, err := s.pciSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get pci data")
	}
	inv.PCIBusDevices = data
	return nil
}

func (s *Svc) setIPMIDevices(inv *inventory.Inventory) error {
	data, err := s.ipmiSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get ipmi data")
	}
	inv.IPMIDevices = data
	return nil
}

func (s *Svc) setNICs(inv *inventory.Inventory) error {
	data, err := s.nicSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get nic data")
	}
	inv.NICs = data
	return nil
}

func (s *Svc) setLLDPFrames(inv *inventory.Inventory) error {
	data, err := s.lldpSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get lldp data")
	}
	inv.LLDPFrames = data
	return nil
}

func (s *Svc) setNDPFrames(inv *inventory.Inventory) error {
	data, err := s.netlinkSvc.GetIPv6NeighbourData()
	if err != nil {
		return errors.Wrap(err, "unable to get ndp data")
	}
	inv.NDPFrames = data
	return nil
}

func (s *Svc) setVirt(inv *inventory.Inventory) error {
	data, err := s.virtSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get virt")
	}
	inv.Virtualization = data
	return nil
}

func (s *Svc) setHost(inv *inventory.Inventory) error {
	hostInfo, err := s.hostSvc.GetData()
	if err != nil {
		return errors.Wrap(err, "unable to get host info")
	}
	inv.Host = hostInfo
	return nil
}

func (s *Svc) setDistro(inv *inventory.Inventory) error {
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
