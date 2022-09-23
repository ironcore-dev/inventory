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

	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks"
	"github.com/onmetal/inventory/internal/benchmarks/output"
)

type Worker interface {
	Do(task Job, result chan<- output.Result)
	ObtainResources(quota string) bool
}

type JobQueue chan Job

type Job struct {
	Bench *conf.Benchmark
}

type workers struct {
	ctx  context.Context
	log  logger.Logger
	Pool *Pool
}

func New(ctx context.Context, l logger.Logger) Worker {
	p := newPool()
	go p.Start()

	return &workers{
		ctx:  ctx,
		log:  l,
		Pool: p,
	}
}

func (wr *workers) Do(task Job, result chan<- output.Result) {
	taskResult := benchmarks.New(task.Bench, wr.log).Start()
	taskResult.BenchmarkName = task.Bench.Name
	taskResult.Name = task.Bench.JSONPathInputSelector
	taskResult.OutputSelector = task.Bench.Output
	result <- taskResult
	wr.Pool.Release <- release{task.Bench.Resources.CPU}
}

func (wr *workers) ObtainResources(quota string) bool {
	return wr.Pool.obtainResources(quota)
}
