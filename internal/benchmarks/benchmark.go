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

package benchmarks

import (
	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks/executor"
	"github.com/onmetal/inventory/internal/benchmarks/output"
	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"
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
