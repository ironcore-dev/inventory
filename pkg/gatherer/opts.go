// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package gatherer

import (
	"github.com/onmetal/inventory/pkg/block"
	"github.com/onmetal/inventory/pkg/cpu"
	"github.com/onmetal/inventory/pkg/distro"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/ipmi"
	"github.com/onmetal/inventory/pkg/lldp"
	"github.com/onmetal/inventory/pkg/mem"
	"github.com/onmetal/inventory/pkg/mlc"
	"github.com/onmetal/inventory/pkg/netlink"
	"github.com/onmetal/inventory/pkg/nic"
	"github.com/onmetal/inventory/pkg/numa"
	"github.com/onmetal/inventory/pkg/pci"
	"github.com/onmetal/inventory/pkg/virt"
)

type Option func(svc *Svc)

func WithDMI(dmiSvc *dmi.Svc) Option {
	return func(svc *Svc) {
		svc.dmiSvc = dmiSvc
	}
}

func WithNUMA(numaSvc *numa.Svc) Option {
	return func(svc *Svc) {
		svc.numaSvc = numaSvc
	}
}

func WithBlocks(blockSvc *block.Svc) Option {
	return func(svc *Svc) {
		svc.blockSvc = blockSvc
	}
}

func WithPCI(pciSvc *pci.Svc) Option {
	return func(svc *Svc) {
		svc.pciSvc = pciSvc
	}
}

func WithCPU(cpuInfoSvc *cpu.InfoSvc) Option {
	return func(svc *Svc) {
		svc.cpuInfoSvc = cpuInfoSvc
	}
}

func WithMem(memInfoSvc *mem.InfoSvc) Option {
	return func(svc *Svc) {
		svc.memInfoSvc = memInfoSvc
	}
}

func WithMLCPerf(mlcPerfSvc *mlc.PerfSvc) Option {
	return func(svc *Svc) {
		svc.mlcPerfSvc = mlcPerfSvc
	}
}

func WithLLDP(lldpSvc *lldp.Svc) Option {
	return func(svc *Svc) {
		svc.lldpSvc = lldpSvc
	}
}

func WithNIC(nicSvc *nic.Svc) Option {
	return func(svc *Svc) {
		svc.nicSvc = nicSvc
	}
}

func WithIPMI(ipmiSvc *ipmi.Svc) Option {
	return func(svc *Svc) {
		svc.ipmiSvc = ipmiSvc
	}
}

func WithNetlink(netlinkSvc *netlink.Svc) Option {
	return func(svc *Svc) {
		svc.netlinkSvc = netlinkSvc
	}
}

func WithVirt(virtSvc *virt.Svc) Option {
	return func(svc *Svc) {
		svc.virtSvc = virtSvc
	}
}

func WithHost(hostSvc *host.Svc) Option {
	return func(svc *Svc) {
		svc.hostSvc = hostSvc
	}
}

func WithDistro(distroSvc *distro.Svc) Option {
	return func(svc *Svc) {
		svc.distroSvc = distroSvc
	}
}
