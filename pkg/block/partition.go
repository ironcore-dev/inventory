// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
