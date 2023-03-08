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

package mlc

import (
	"github.com/pkg/errors"
)

const (
	CMemPerfLocalMemBWKey       = "LocalMemBW"
	CMemPerfRemoteMemBWKey      = "RemoteMemBW"
	CMemPerfLocalMemLatencyKey  = "LocalMemLatency"
	CMemPerfRemoteMemLatencyKey = "RemoteMemLatency"
)

type Perf struct {
	LocalMemBW       float64
	RemoteMemBW      float64
	LocalMemLatency  float64
	RemoteMemLatency float64
}

func (memperf *Perf) setField(key string, val float64) error {
	switch key {
	case CMemPerfLocalMemBWKey:
		memperf.LocalMemBW = val
	case CMemPerfRemoteMemBWKey:
		memperf.RemoteMemBW = val
	case CMemPerfLocalMemLatencyKey:
		memperf.LocalMemLatency = val
	case CMemPerfRemoteMemLatencyKey:
		memperf.RemoteMemLatency = val
	default:
		return errors.Errorf("unknown key %s from meminfo", key)
	}
	return nil
}
