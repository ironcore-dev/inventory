// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
