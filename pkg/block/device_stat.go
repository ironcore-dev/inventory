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

package block

import (
	"github.com/pkg/errors"
)

// https://www.kernel.org/doc/Documentation/block/stat.txt

type DeviceStat struct {
	ReadIOs        uint64
	ReadMerges     uint64
	ReadSectors    uint64
	ReadTicks      uint64
	WriteIOs       uint64
	WriteMerges    uint64
	WriteSectors   uint64
	WriteTicks     uint64
	InFlight       uint64
	IOTicks        uint64
	TimeInQueue    uint64
	DiscardIOs     uint64
	DiscardMerges  uint64
	DiscardSectors uint64
	DiscardTicks   uint64
	FlushIOs       uint64
	FlushTicks     uint64
}

func (s *DeviceStat) setByIndex(idx int, val uint64) error {
	switch idx {
	case 0:
		s.ReadIOs = val
	case 1:
		s.ReadMerges = val
	case 2:
		s.ReadSectors = val
	case 3:
		s.ReadTicks = val
	case 4:
		s.WriteIOs = val
	case 5:
		s.WriteMerges = val
	case 6:
		s.WriteSectors = val
	case 7:
		s.WriteTicks = val
	case 8:
		s.InFlight = val
	case 9:
		s.IOTicks = val
	case 10:
		s.TimeInQueue = val
	case 11:
		s.DiscardIOs = val
	case 12:
		s.DiscardMerges = val
	case 13:
		s.DiscardSectors = val
	case 14:
		s.DiscardTicks = val
	case 15:
		s.FlushIOs = val
	case 16:
		s.FlushTicks = val
	default:
		return errors.Errorf("unexpected index %d", idx)
	}
	return nil
}
