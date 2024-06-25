// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
