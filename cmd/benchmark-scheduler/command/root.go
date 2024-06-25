// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"context"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/dispatcher"
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
