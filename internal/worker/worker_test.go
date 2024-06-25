// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"os/user"
	"testing"

	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks/output"
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
