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
	"context"
	"os/user"
	"testing"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks/output"
	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"
)

func TestNew(t *testing.T) {
	l := logger.New()
	n := New(context.Background(), l)
	if !n.ObtainResources(quota) {
		t.Log("can't obtain resources from pool")
		t.Fail()
	}
}

func TestWorkerDo(t *testing.T) {
	u, err := user.Current()
	if err != nil || u.Name != "root" {
		t.Skipf("can't get user info or user is not a root: %s", err)
	}

	l := logger.New()
	n := New(context.Background(), l)

	job := Job{
		Bench: &conf.Benchmark{
			Name: "test", Application: "echo", JSONPathInputSelector: "echo",
			Args:      []string{"hello", "go"},
			Resources: &conf.Resources{CPUS: conf.CPUS{CPU: "100m"}},
		},
	}
	jobResult := make(chan output.Result, 1)
	n.Do(job, jobResult)

	r := <-jobResult

	if r.Name != "echo" {
		t.Log("field mismatch", "was", "echo", "got", r.Name)
		t.Fail()
	}
	if r.Error != nil {
		t.Log(r.Error)
		t.Fail()
	}
}
