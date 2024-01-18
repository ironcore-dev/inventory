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
	"errors"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/google/uuid"
	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/stretchr/testify/assert"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks/output"
	"github.com/onmetal/inventory/internal/provider/fake"
)

var errFilePathDefinition = errors.New("can't define test location")

func TestDo(t *testing.T) {
	a := assert.New(t)
	l := logger.New()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Log(errFilePathDefinition)
		t.Fail()
	}
	dir := path.Join(path.Dir(filename), "testdata/fio.json")

	data, err := os.ReadFile(dir)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fakeClient := fake.New()
	m := Machine{
		uuid:      "test",
		log:       l,
		results:   getResults(data),
		Client:    fakeClient,
		resultMap: make(map[string]metalv1alpha4.Benchmarks),
	}
	if err := m.Do(); !a.Nil(err) {
		t.Log(err)
		t.Fail()
	}
}

func getResults(data []byte) []output.Result {
	return []output.Result{
		{
			OutputSelector: "jobs.*.read.io_bytes",
			Message:        data,
			UUID:           uuid.New(),
		},
	}
}
