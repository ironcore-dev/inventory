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
	"context"
	"sync"
	"time"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks/output"
	"github.com/onmetal/inventory/internal/worker"
)

const (
	timeOutSecond   = 600 * time.Second
	tickTimerSecond = 60 * time.Second
)

// Dispatcher is the interface that groups the basic methods for interacting.
type Dispatcher interface {
	Start()
	AddJob(job worker.Job)
	RequeueJob(task worker.Job)
	JobResult() output.JobResultQueue
}

type disp struct {
	ctx       context.Context
	rwm       *sync.RWMutex
	job       worker.JobQueue
	jobResult output.JobResultQueue
	Requeue   worker.JobQueue
	log       logger.Logger
	worker    worker.Worker
}

func NewWithSize(queueSize int, l logger.Logger) Dispatcher {
	ctx := context.Background()
	return &disp{
		ctx:       ctx,
		rwm:       new(sync.RWMutex),
		job:       make(worker.JobQueue, queueSize),
		jobResult: make(output.JobResultQueue, queueSize),
		Requeue:   make(worker.JobQueue, queueSize),
		log:       l,
	}
}

func (d *disp) Start() {
	d.log.Info("dispatcher started")

	d.worker = worker.New(d.ctx, d.log)
	for {
		select {
		case job := <-d.Requeue:
			go d.run(job)
		case job := <-d.job:
			go d.run(job)
		}
	}
}

func (d *disp) AddJob(job worker.Job) {
	d.job <- job
}

func (d *disp) RequeueJob(task worker.Job) {
	timeout := time.After(timeOutSecond)
	tick := time.NewTicker(tickTimerSecond)
	for {
		select {
		case <-timeout:
			return
		case <-tick.C:
			d.log.Info("trying to requeue task", "name", task.Bench.Name, "args", task.Bench.Args)
			d.Requeue <- task
			return
		}
	}
}

func (d *disp) JobResult() output.JobResultQueue {
	for {
		return d.jobResult
	}
}

func (d *disp) run(job worker.Job) {
	if !d.worker.ObtainResources(job.Bench.Resources.CPU) {
		d.log.Info("not enough quota to start new job",
			"name", job.Bench.Name, "command", job.Bench.Args)
		go d.RequeueJob(job)
		return
	}
	go d.worker.Do(job, d.jobResult)
}