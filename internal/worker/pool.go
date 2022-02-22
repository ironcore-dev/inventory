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

package worker

import (
	"runtime"
	"sync/atomic"

	"github.com/onmetal/inventory/internal/strconverter"
)

const coresMultiplier int64 = 100000

type Pool struct {
	cpus           map[string]string
	Obtain         chan string
	Release        chan release
	size, capacity int64
}

type release struct {
	quota string
}

func newPool() *Pool {
	osCoresCount := runtime.NumCPU()
	cpuSet := make(map[string]string, osCoresCount)

	return &Pool{
		size:     strconverter.ServerCPUFullCapacity,
		capacity: strconverter.ServerCPUFullCapacity,
		Release:  make(chan release, osCoresCount),
		Obtain:   make(chan string, osCoresCount),
		cpus:     cpuSet,
	}
}

func (c *Pool) Start() {
	for {
		r := <-c.Release
		c.updateSize(c.getSize() + strconverter.QuotaToInt(r.quota))
	}
}

func (c *Pool) obtainResources(quota string) bool {
	if quota == "all" {
		if (c.getSize() - c.capacity) < 0 {
			return false
		}
		c.updateSize(c.getSize() - c.capacity)
		return true
	}
	q := strconverter.QuotaToInt(quota)
	if (c.getSize() - q) < 0 {
		return false
	}
	c.updateSize(c.getSize() - q)
	return true
}

func (c *Pool) getSize() int64 {
	return atomic.LoadInt64(&c.size)
}

func (c *Pool) updateSize(newSize int64) {
	atomic.SwapInt64(&c.size, newSize)
}
