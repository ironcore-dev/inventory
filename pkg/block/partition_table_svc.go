// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package block

import (
	"path"
	"strings"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/partition/gpt"
	"github.com/diskfs/go-diskfs/partition/mbr"
	"github.com/pkg/errors"
)

const (
	CDevBasePath = "/dev"

	CGPTPartitionTableType = "GPT"
	CMBRPartitionTableType = "MBR"
)

type PartitionTableSvc struct {
	devPath string
}

func NewPartitionTableSvc(basePath string) *PartitionTableSvc {
	return &PartitionTableSvc{
		devPath: path.Join(basePath, CDevBasePath),
	}
}

func (s *PartitionTableSvc) GetPartitionTable(devName string) (*PartitionTable, error) {
	devPath := path.Join(s.devPath, devName)
	disk, err := diskfs.Open(devPath, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open device %s", devPath)
	}

	recTable, err := disk.GetPartitionTable()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get disk partition table")
	}

	tableType := strings.ToUpper(recTable.Type())

	var table *PartitionTable
	switch tableType {
	case CGPTPartitionTableType:
		realTable := recTable.(*gpt.Table)
		table = &PartitionTable{
			Type:       CGPTPartitionTableType,
			Partitions: NewPartitionsFromGPT(realTable.Partitions),
		}
	case CMBRPartitionTableType:
		realTable := recTable.(*mbr.Table)
		table = &PartitionTable{
			Type:       CMBRPartitionTableType,
			Partitions: NewPartitionsFromMBR(realTable.Partitions),
		}
	default:
		return nil, errors.Errorf("unsupported partition table type %s", tableType)
	}

	return table, nil
}
