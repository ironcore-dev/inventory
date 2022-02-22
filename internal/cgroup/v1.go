// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package cgroup

import (
	"fmt"

	"github.com/containerd/cgroups"
	"github.com/onmetal/inventory/internal/strconverter"
	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"
	"github.com/opencontainers/runtime-spec/specs-go"
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
