// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package cgroup

import (
	"fmt"

	"github.com/containerd/cgroups"
	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"
	"github.com/opencontainers/runtime-spec/specs-go"

	"github.com/onmetal/inventory/internal/strconverter"
)

type v1 struct {
	cgroups.Cgroup
}

func newV1(name, mountPoint string, resources *conf.Resources) (*v1, error) {
	cGroupName := fmt.Sprintf("%s/%s", mountPoint, name)
	quota := strconverter.QuotaToInt(resources.CPU)

	period := resources.Period
	if period == 0 {
		period = linuxKernelDefaultPeriod
	}
	weight := resources.Shares
	if weight == 0 {
		weight = linuxKernelDefaultWeight
	}

	controller, err := cgroups.New(cgroups.V1, cgroups.StaticPath(cGroupName), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &weight,
			Quota:  &quota,
			Period: &period,
			Cpus:   resources.CPUSet,
		},
		Memory: &specs.LinuxMemory{
			Swap:  resources.Swap,
			Limit: resources.Max,
		},
	})
	return &v1{controller}, err
}

func (c *v1) Add(pid int) error {
	return c.Cgroup.Add(cgroups.Process{Pid: pid})
}

func (c *v1) Delete() error {
	return c.Cgroup.Delete()
}
