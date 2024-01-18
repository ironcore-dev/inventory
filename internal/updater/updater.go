// Copyright 2023 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package updater

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/Jeffail/gabs/v2"
	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks/output"
	"github.com/onmetal/inventory/internal/provider"
)

const (
	minimalSliceLen = 2
)

type Machine struct {
	provider.Client

	resultMap map[string]metalv1alpha4.Benchmarks
	uuid      string
	wg        *sync.WaitGroup
	log       logger.Logger
	results   []output.Result
}

type Client interface {
	Do() error
}

func New(machineUUID string, results []output.Result, c provider.Client, l logger.Logger) (*Machine, error) {
	return &Machine{
		Client:    c,
		wg:        new(sync.WaitGroup),
		uuid:      machineUUID,
		results:   results,
		resultMap: make(map[string]metalv1alpha4.Benchmarks, len(results)),
		log:       l,
	}, nil
}

func (m *Machine) Do() error {
	for r := range m.results {
		if m.results[r].Error != nil {
			m.log.Info("job exited with error", "error", m.results[r].Error)
			continue
		}
		name := m.results[r].Name
		if name == "" {
			name = m.results[r].BenchmarkName
		}
		switch {
		case strings.Contains(m.results[r].OutputSelector, "text"):
			value := m.parseText(&m.results[r])
			benches, ok := m.resultMap[name]
			if !ok {
				m.resultMap[name] = metalv1alpha4.Benchmarks{{Name: m.results[r].BenchmarkName, Value: value}}
				continue
			}
			m.resultMap[name] = append(benches, metalv1alpha4.BenchmarkResult{
				Name: m.results[r].BenchmarkName, Value: value,
			})
		default:
			value := m.parseJSON(&m.results[r])
			benches, ok := m.resultMap[name]
			if !ok {
				m.resultMap[name] = metalv1alpha4.Benchmarks{{Name: m.results[r].BenchmarkName, Value: value}}
			}
			m.resultMap[name] = append(benches, metalv1alpha4.BenchmarkResult{Name: m.results[r].BenchmarkName, Value: value})
		}
	}
	patch := metalv1alpha4.Benchmark{Spec: metalv1alpha4.BenchmarkSpec{Benchmarks: m.resultMap}}
	body, err := json.Marshal(patch)
	if err != nil {
		return err
	}
	m.log.Info("machine benchmark updating", "name", m.uuid)
	return m.Client.Patch(m.uuid, body)
}

func (m *Machine) parseText(res *output.Result) uint64 {
	filter := strings.Split(res.OutputSelector, ":")
	if len(filter) != minimalSliceLen {
		m.log.Info("output split failed", "name", res.BenchmarkName, "output", res.OutputSelector)
		return 0
	}
	splittedMessage := strings.FieldsFunc(string(res.Message), split)
	if len(splittedMessage) > minimalSliceLen {
		m.log.Info("output split failed", "name", res.BenchmarkName, "output", string(res.Message))
		return 0
	}
	return getValueFromText(splittedMessage, filter[1])
}

func (m *Machine) parseJSON(run *output.Result) uint64 {
	parsed, err := gabs.ParseJSON(run.Message)
	if err != nil {
		m.log.Info("failed to parse json", "error", err)
		return 0
	}
	return getValueFromJSON(parsed.Path(run.OutputSelector).Data())
}

func split(r rune) bool {
	return r == ' ' || r == '\n'
}

func getValueFromText(splittedMessage []string, filter string) uint64 {
	for s := range splittedMessage {
		if !strings.Contains(splittedMessage[s], filter) {
			continue
		}
		value, err := strconv.Atoi(splittedMessage[s+1])
		if err != nil {
			return 0
		}
		return uint64(value)
	}
	return 0
}

func getValueFromJSON(v interface{}) uint64 {
	switch t := v.(type) {
	case []interface{}:
		n, ok := t[0].(float64)
		if !ok {
			return 0
		}
		return uint64(n)
	case interface{}:
		n, ok := t.(float64)
		if !ok {
			return 0
		}
		return uint64(n)
	default:
		return 0
	}
}
