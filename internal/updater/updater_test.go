// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
