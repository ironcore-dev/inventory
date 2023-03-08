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

type Device struct {
	path              string
	Name              string
	Type              string
	Rotational        bool
	Removable         bool
	ReadOnly          bool
	Vendor            string
	Model             string
	Serial            string
	WWID              string
	FirmwareRevision  string
	State             string
	PhysicalBlockSize uint64
	LogicalBlockSize  uint64
	HWSectorSize      uint64
	Size              uint64
	NUMANodeID        uint64
	PartitionTable    *PartitionTable
	Stat              *DeviceStat
}
