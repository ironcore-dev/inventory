// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package block

type PartitionTableType string

type PartitionTable struct {
	Type       PartitionTableType
	Partitions []Partition
}
