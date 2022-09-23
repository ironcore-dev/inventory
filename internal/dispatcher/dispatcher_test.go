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

package dispatcher

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"
	"testing"
	"time"

	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/worker"
)

func TestNewWithSize(t *testing.T) {
	u, err := user.Current()
	if err != nil || u.Name != "root" {
		t.Skipf("can't get user info or user is not a root: %s", err)
	}

	l := logger.New()
	n := NewWithSize(runtime.NumCPU(), l)

	go n.Start()
	time.Sleep(100 * time.Millisecond)

	cores := runtime.NumCPU()
	job := []worker.Job{
		{
			Bench: &conf.Benchmark{
				Name: "test2", Application: "echo", Args: []string{"hello", "test2"},
				Resources: &conf.Resources{CPUS: conf.CPUS{CPU: fmt.Sprintf("%d", cores)}},
			},
		},
	}
	for j := range job {
		n.AddJob(job[j])
	}
	time.Sleep(100 * time.Millisecond)
	for i := 1; i < len(job); i++ {
		tr := <-n.JobResult()
		if !strings.Contains(string(tr.Message), "hello test") {
			t.Log("output is modified", "got", string(tr.Message), "must contain: `hello test`")
			t.Fail()
		}
	}
}
