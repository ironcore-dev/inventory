package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	"github.com/stretchr/testify/assert"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
)

func TestPatch(t *testing.T) {
	a := assert.New(t)
	l := logger.New()
	ctx := context.Background()

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		contains(t, req.URL.String(), "/apis/v1alpha3/benchmark/onmetal/test")
		// Send response to be tested
		rw.Header().Set("Content-Type", "application/json")
		//nolint:errcheck
		rw.Write([]byte(`OK`))
	}))

	server.Start()
	time.Sleep(1 * time.Second)
	defer server.Close()

	h := httpClient{
		Client: server.Client(), gateway: fmt.Sprintf("http://%s", server.Listener.Addr().String()),
		namespace: "onmetal",
		ctx:       ctx, log: l,
	}

	mm := make(map[string]benchv1alpha3.Benchmarks)
	mm["test"] = []benchv1alpha3.Benchmark{{Name: "disk-test", Value: 123}}
	spec := benchv1alpha3.Machine{Spec: benchv1alpha3.MachineSpec{Benchmarks: mm}}

	patch, err := json.Marshal(spec)
	a.Nil(err, "json serialization failed")

	a.Nil(h.Patch("test", patch), "object patch failed")
}

func contains(tb testing.TB, exp, act interface{}) {
	if !strings.Contains(exp.(string), act.(string)) {
		tb.FailNow()
	}
}
