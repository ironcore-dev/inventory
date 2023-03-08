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
