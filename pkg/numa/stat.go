// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package numa

import (
	"github.com/pkg/errors"
)

const (
	CNodeNumaHitKey       = "numa_hit"
	CNodeNumaMissKey      = "numa_miss"
	CNodeNumaForeignKey   = "numa_foreign"
	CNodeInterleaveHitKey = "interleave_hit"
	CNodeLocalNodeKey     = "local_node"
	CNodeOtherNodeKey     = "other_node"
)

type Stat struct {
	NumaHit       uint64
	NumaMiss      uint64
	NumaForeign   uint64
	InterleaveHit uint64
	LocalNode     uint64
	OtherNode     uint64
}

func (stat *Stat) setField(key string, val uint64) error {
	switch key {
	case CNodeNumaHitKey:
		stat.NumaHit = val
	case CNodeNumaMissKey:
		stat.NumaMiss = val
	case CNodeNumaForeignKey:
		stat.NumaForeign = val
	case CNodeInterleaveHitKey:
		stat.InterleaveHit = val
	case CNodeLocalNodeKey:
		stat.LocalNode = val
	case CNodeOtherNodeKey:
		stat.OtherNode = val
	default:
		return errors.Errorf("unknown key %s from meminfo", key)
	}
	return nil
}
