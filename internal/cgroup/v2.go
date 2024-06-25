// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package cgroup

import (
	"fmt"

	cgroupsv2 "github.com/containerd/cgroups/v2"
	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"

	"github.com/onmetal/inventory/internal/strconverter"
)

const (
	linuxKernelDefaultPeriod = 100000
	linuxKernelDefaultWeight = 100
)

type v2 struct {
	*cgroupsv2.Manager
}

func newV2(name, mountPoint string, resources *conf.Resources) (*v2, error) {
	cGroupName := fmt.Sprintf("/%s", name)
	quota := strconverter.QuotaToInt(resources.CPU)
	period := resources.Period
	if period == 0 {
		period = linuxKernelDefaultPeriod
	}
	weight := resources.Shares
	if weight == 0 {
		weight = linuxKernelDefaultWeight
	}
	m, err := cgroupsv2.NewManager(mountPoint, cGroupName, &cgroupsv2.Resources{
		CPU: &cgroupsv2.CPU{
			Weight: &weight,
			Max:    cgroupsv2.NewCPUMax(&quota, &period),
			Cpus:   resources.CPUSet,
		},
		Memory: &cgroupsv2.Memory{
			Swap: resources.Swap,
			Max:  resources.Max,
			Low:  resources.Low,
			High: resources.High,
		},
	})
	return &v2{m}, err
}

func (c *v2) Add(pid int) error {
	return c.Manager.AddProc(uint64(pid))
}

func (c *v2) Delete() error {
	return c.Manager.Delete()
}
