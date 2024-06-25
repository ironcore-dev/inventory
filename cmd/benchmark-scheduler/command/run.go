// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks/output"
	"github.com/onmetal/inventory/internal/config"
	"github.com/onmetal/inventory/internal/provider"
	"github.com/onmetal/inventory/internal/smbiosinfo"
	"github.com/onmetal/inventory/internal/updater"
	"github.com/onmetal/inventory/internal/worker"
)

func (b *BenchOpts) newRun() *cli.Command {
	return &cli.Command{
		Name:    "run",
		Aliases: []string{"start", "do"},
		Usage:   "run benchmark jobs",
		Action:  b.run,
		Flags:   checkFlags(),
	}
}

func (b *BenchOpts) run(cliCtx *cli.Context) error {
	b.log.Info("program started")

	machineUUID, uuidErr := b.getMachineUUID()
	if uuidErr != nil {
		b.log.Info("can't get UUID from machine", "error", uuidErr)
		return uuidErr
	}
	b.machineUUID = machineUUID

	c, err := provider.New(context.Background(), b.log, cliCtx)
	if err != nil {
		return err
	}

	go b.dispatcher.Start()

	conf := config.New(machineUUID, cliCtx, c, b.log)

	for task := range conf.Benchmarks {
		b.log.Info("task added", "name", conf.Benchmarks[task].Name)
		j := worker.Job{Bench: &conf.Benchmarks[task]}
		b.dispatcher.AddJob(j)
	}
	r := b.waitForResults(len(conf.Benchmarks))

	return b.update(r, c, b.log)
}

func (b *BenchOpts) waitForResults(tasks int) []output.Result {
	b.log.Info("waiting for result")
	r := make([]output.Result, 0, tasks)
	for i := 1; i <= tasks; i++ {
		r = append(r, <-b.dispatcher.JobResult())
	}
	return r
}

func (b *BenchOpts) update(r []output.Result, c provider.Client, l logger.Logger) error {
	u, err := updater.New(b.machineUUID, r, c, l)
	if err != nil {
		return err
	}
	return u.Do()
}

func (b *BenchOpts) getMachineUUID() (string, error) {
	sm, err := smbiosinfo.New(b.log)
	if err != nil {
		return "", err
	}
	defer func(sm smbiosinfo.SystemManager) {
		if err := sm.Close(); err != nil {
			b.log.Info("can't close stream properly", "error", err)
		}
	}(sm)
	return sm.GetUUID()
}
