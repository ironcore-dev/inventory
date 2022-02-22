// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package worker

import (
	"runtime"
	"testing"
	"time"
)

const quota = "100m"

func TestResourceRelease(t *testing.T) {
	osCoresCount := runtime.NumCPU()
	size := int64(osCoresCount) * coresMultiplier

	p := newPool()
	go p.Start()
	// wait for start
	time.Sleep(100 * time.Millisecond)

	if !p.obtainResources(quota) {
		t.Log("can't obtain resources from pool")
		t.Fail()
	}

	if p.getSize() != size-10000 {
		t.Log("invalid output value", p.getSize(), "got", size-10000)
		t.Fail()
	}

	p.Release <- release{quota: quota}
	// wait for release
	time.Sleep(100 * time.Millisecond)

	if p.getSize() != size {
		t.Log("invalid output value", size, "got", p.getSize())
		t.Fail()
	}
}

func TestNoResourceAvailable(t *testing.T) {
	quota := "100m"

	p := newPool()
	p.updateSize(0)
	if p.obtainResources(quota) {
		t.Log("there are should be no resources available for obtain")
		t.Fail()
	}
}
