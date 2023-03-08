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
	"strconv"

	"github.com/diskfs/go-diskfs/partition/gpt"
	"github.com/diskfs/go-diskfs/partition/mbr"
)

type Partition struct {
	ID          string
	Name        string
	Size        uint64
	StartSector uint64
	EndSector   uint64
}

func NewPartitionsFromMBR(recPartitions []*mbr.Partition) []Partition {
	var partitions []Partition
	for i, recPartition := range recPartitions {
		part := NewPartitionFromMBR(strconv.Itoa(i), recPartition)
		partitions = append(partitions, *part)
	}
	return partitions
}

func NewPartitionFromMBR(id string, partition *mbr.Partition) *Partition {
	return &Partition{
		ID:          id,
		Size:        uint64(partition.GetSize()),
		StartSector: uint64(partition.StartSector),
		EndSector:   uint64(partition.EndSector),
	}
}

func NewPartitionsFromGPT(recPartitions []*gpt.Partition) []Partition {
	var partitions []Partition
	for _, recPartition := range recPartitions {
		if recPartition.Type == gpt.Unused {
			continue
		}
		part := NewPartitionFromGPT(recPartition)
		partitions = append(partitions, *part)
	}
	return partitions
}

func NewPartitionFromGPT(partition *gpt.Partition) *Partition {
	return &Partition{
		ID:          partition.GUID,
		Name:        partition.Name,
		Size:        uint64(partition.GetSize()),
		StartSector: partition.Start,
		EndSector:   partition.End,
	}
}
