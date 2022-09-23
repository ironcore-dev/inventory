package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/block"
	"github.com/onmetal/inventory/pkg/cpu"
	"github.com/onmetal/inventory/pkg/crd"
	"github.com/onmetal/inventory/pkg/distro"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/flags"
	"github.com/onmetal/inventory/pkg/gatherer"
	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/ipmi"
	"github.com/onmetal/inventory/pkg/lldp"
	"github.com/onmetal/inventory/pkg/lldp/frame"
	"github.com/onmetal/inventory/pkg/mem"
	"github.com/onmetal/inventory/pkg/netlink"
	"github.com/onmetal/inventory/pkg/nic"
	"github.com/onmetal/inventory/pkg/numa"
	"github.com/onmetal/inventory/pkg/pci"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/redis"
	"github.com/onmetal/inventory/pkg/virt"
)

type InventoryApp struct {
	printer       *printer.Svc
	gathererSvc   *gatherer.Svc
	crdBuilderSvc *crd.BuilderSvc
	crdSaverSvc   crd.SaverSvc
	crdSaverPatch bool
}

func NewInventoryApp() (*InventoryApp, int) {
	f, err := flags.NewInventoryFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, CErrRetCode
	}
	p := printer.NewSvc(f.Verbose)

	crdBuilderSvc := crd.NewBuilderSvc(p)

	var crdSvcConstructor func() (crd.SaverSvc, error)
	if f.Gateway != "" {
		crdSvcConstructor = func() (crd.SaverSvc, error) {
			return crd.NewGatewaySaverSvc(f.Gateway, f.KubeNamespace, f.Timeout)
		}
	} else {
		crdSvcConstructor = func() (crd.SaverSvc, error) {
			return crd.NewKubeAPISaverSvc(f.Kubeconfig, f.KubeNamespace)
		}
	}

	crdSaverSvc, err := crdSvcConstructor()
	if err != nil {
		p.Err(errors.Wrapf(err, "unable to create k8s resorce saver svc"))
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

	redisSvc := redis.NewRedisSvc(f.Root, f.RedisUser, f.RedisPassword)
	lldpFrameInfoSvc := frame.NewFrameSvc(p)
	lldpSvc := lldp.NewSvc(p, lldpFrameInfoSvc, hostSvc, redisSvc, f.Root)

	nicDevSvc := nic.NewDeviceSvc(p)
	nicSvc := nic.NewSvc(p, nicDevSvc, hostSvc, redisSvc, f.Root)

	ipmiDevInfoSvc := ipmi.NewDeviceSvc(p)
	ipmiSvc := ipmi.NewSvc(p, ipmiDevInfoSvc, f.Root)

	nlSvc := netlink.NewSvc(p, f.Root)

	virtSvc := virt.NewSvc(dmiSvc, cpuInfoSvc, f.Root)

	distroSvc := distro.NewSvc(p, hostSvc, f.Root)

	opts := []gatherer.Option{
		gatherer.WithDMI(dmiSvc),
		gatherer.WithNUMA(numaSvc),
		gatherer.WithBlocks(blockSvc),
		gatherer.WithPCI(pciSvc),
		gatherer.WithCPU(cpuInfoSvc),
		gatherer.WithMem(memInfoSvc),
		gatherer.WithLLDP(lldpSvc),
		gatherer.WithNIC(nicSvc),
		gatherer.WithIPMI(ipmiSvc),
		gatherer.WithNetlink(nlSvc),
		gatherer.WithVirt(virtSvc),
		gatherer.WithHost(hostSvc),
		gatherer.WithDistro(distroSvc),
	}

	gathererSvc := gatherer.NewSvc(p, opts...)

	return &InventoryApp{
		printer:       p,
		gathererSvc:   gathererSvc,
		crdBuilderSvc: crdBuilderSvc,
		crdSaverSvc:   crdSaverSvc,
		crdSaverPatch: f.Patch,
	}, 0
}

func (s *InventoryApp) Run() int {
	inv := s.gathererSvc.Gather()

	cr, err := s.crdBuilderSvc.Build(inv)
	if err != nil {
		s.printer.Err(errors.Wrap(err, "unable to build inventory resource"))
		return CErrRetCode
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

	if s.crdSaverPatch {
		err = s.crdSaverSvc.Patch(inv.DMI.SystemInformation.UUID, cr)
	} else {
		err = s.crdSaverSvc.Save(cr)
	}
	if err != nil {
		s.printer.Err(errors.Wrap(err, "unable to save inventory resource"))
		return CErrRetCode
	}

	return COKRetCode
}
