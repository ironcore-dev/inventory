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

package netlink

import "github.com/vishvananda/netlink"

const (
	CNeighbourNoneCacheState       = "None"
	CNeighbourIncompleteCacheState = "Incomplete"
	CNeighbourReachableCacheState  = "Reachable"
	CNeighbourStaleCacheState      = "Stale"
	CNeighbourDelayCacheState      = "Delay"
	CNeighbourProbeCacheState      = "Probe"
	CNeighbourFailedCacheState     = "Failed"
	CNeighbourNoARPCacheState      = "No ARP"
	CNeighbourPermanentCacheState  = "Permanent"
)

type NeighbourCacheState string

var CNeighbourCacheStates = map[int]NeighbourCacheState{
	0x00: CNeighbourNoneCacheState,
	0x01: CNeighbourIncompleteCacheState,
	0x02: CNeighbourReachableCacheState,
	0x04: CNeighbourStaleCacheState,
	0x08: CNeighbourDelayCacheState,
	0x10: CNeighbourProbeCacheState,
	0x20: CNeighbourFailedCacheState,
	0x40: CNeighbourNoARPCacheState,
	0x80: CNeighbourPermanentCacheState,
}

type IPv6Neighbour struct {
	DeviceIndex int
	DeviceName  string
	IP          string
	MACAddress  string
	State       NeighbourCacheState
}

func NewIPv6Neighbour(idx int, name string, n *netlink.Neigh) *IPv6Neighbour {
	neighbour := &IPv6Neighbour{
		DeviceIndex: idx,
		DeviceName:  name,
		IP:          n.IP.String(),
		MACAddress:  n.HardwareAddr.String(),
		State:       CNeighbourCacheStates[n.State],
	}
	return neighbour
}
