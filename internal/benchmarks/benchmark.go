// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks/executor"
	"github.com/onmetal/inventory/internal/benchmarks/output"
)

type Benchmarker interface {
	Start() output.Result
}

func New(b *conf.Benchmark, l logger.Logger) Benchmarker {
	return &executor.Task{
		Benchmark: b,
		Log:       l,
	}
}
