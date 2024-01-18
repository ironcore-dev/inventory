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

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
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

	mm := make(map[string]metalv1alpha4.Benchmarks)
	mm["test"] = []metalv1alpha4.BenchmarkResult{{Name: "disk-test", Value: 123}}
	spec := metalv1alpha4.Benchmark{Spec: metalv1alpha4.BenchmarkSpec{Benchmarks: mm}}

	patch, err := json.Marshal(spec)
	a.Nil(err, "json serialization failed")

	a.Nil(h.Patch("test", patch), "object patch failed")
}

func contains(tb testing.TB, exp, act interface{}) {
	if !strings.Contains(exp.(string), act.(string)) {
		tb.FailNow()
	}
}
