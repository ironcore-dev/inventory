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

package command

import (
	"context"
	"runtime"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/dispatcher"
	"github.com/urfave/cli/v2"
)

type BenchOpts struct {
	ctx         context.Context
	log         logger.Logger
	dispatcher  dispatcher.Dispatcher
	machineUUID string
}

func NewRoot(version string) cli.App {
	l := logger.New()
	b := newBenchOptions(l)
	return cli.App{
		Name:    "bench-scheduler",
		Usage:   "Start benchmarks in a scheduler way",
		Version: version,
		Commands: []*cli.Command{
			b.newRun(),
		},
	}
}

func newBenchOptions(l logger.Logger) BenchOpts {
	return BenchOpts{
		ctx:        context.Background(),
		dispatcher: dispatcher.NewWithSize(runtime.NumCPU(), l),
		log:        l,
	}
}
