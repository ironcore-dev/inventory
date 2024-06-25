// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package nic

const (
	CNICFlagsMaxFlagCount = 19
)

type Flags struct {
	Up                         bool
	Broadcast                  bool
	Debug                      bool
	Loopback                   bool
	PointToPoint               bool
	NoTrailers                 bool
	Running                    bool
	NoARP                      bool
	Promiscuous                bool
	ReceiveAllMulticastPackets bool
	Master                     bool
	Slave                      bool
	Multicast                  bool
	IfmapSelection             bool
	AutomediaSelection         bool
	DynamicAddress             bool
	LowerUp                    bool
	Dormant                    bool
	Echo                       bool
}

func NewFlags(flagsNum uint32) *Flags {
	flags := &Flags{}

	for i := 0; i < CNICFlagsMaxFlagCount; i++ {
		idx := uint32(1 << i)
		val := flagsNum & idx
		flags.setByIndex(i, val != 0)
	}

	return flags
}

func (s *Flags) setByIndex(idx int, val bool) {
	switch idx {
	case 0:
		s.Up = val
	case 1:
		s.Broadcast = val
	case 2:
		s.Debug = val
	case 3:
		s.Loopback = val
	case 4:
		s.PointToPoint = val
	case 5:
		s.NoTrailers = val
	case 6:
		s.Running = val
	case 7:
		s.NoARP = val
	case 8:
		s.Promiscuous = val
	case 9:
		s.ReceiveAllMulticastPackets = val
	case 10:
		s.Master = val
	case 11:
		s.Slave = val
	case 12:
		s.Multicast = val
	case 13:
		s.IfmapSelection = val
	case 14:
		s.AutomediaSelection = val
	case 15:
		s.DynamicAddress = val
	case 16:
		s.LowerUp = val
	case 17:
		s.Dormant = val
	case 18:
		s.Echo = val
	}
}
