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
