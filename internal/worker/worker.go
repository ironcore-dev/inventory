// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
