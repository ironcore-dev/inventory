package updater

import (
	"errors"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/google/uuid"
	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
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
		resultMap: make(map[string]benchv1alpha3.Benchmarks),
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
