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

package executor

import (
	"bytes"
	"os/exec"

	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/onmetal/inventory/internal/benchmarks/output"
	"github.com/onmetal/inventory/internal/cgroup"
	bencherr "github.com/onmetal/inventory/internal/errors"
	"github.com/onmetal/inventory/internal/strconverter"
)

type Task struct {
	*conf.Benchmark

	Log logger.Logger
}

func (t *Task) Start() output.Result {
	t.Log = t.Log.WithValues("name", t.Name, "args", t.Args)
	t.Log.Info("starting benchmark")

	cGroupUUID := strconverter.RandomString(10)
	control, cgroupErr := cgroup.NewWithLimits(cGroupUUID, t.Resources)
	if cgroupErr != nil {
		t.Log.Info("can't create CGroup", "error", cgroupErr)
		return output.Result{Error: bencherr.CGroupError(cgroupErr)}
	}
	defer func(control cgroup.Manager) {
		if delErr := control.Delete(); delErr != nil {
			t.Log.Info("can't delete CGroup", "error", delErr)
		}
	}(control)

	return t.startWithCGroup(control)
}

func (t *Task) startWithCGroup(control cgroup.Manager) output.Result {
	cmd := exec.Command(t.Application, t.Args...) //nolint:gosec

	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer

	if startErr := cmd.Start(); startErr != nil {
		t.Log.Info("can't start program", "error", startErr)
		return output.Result{Message: nil, Error: bencherr.ExecError(startErr.Error())}
	}

	if cgroupAddErr := control.Add(cmd.Process.Pid); cgroupAddErr != nil {
		t.Log.Info("can't add pid to CGroup", "error", cgroupAddErr)
		return output.Result{Message: outputBuffer.Bytes(), Error: bencherr.CGroupError(cgroupAddErr)}
	}

	if waitErr := cmd.Wait(); waitErr != nil {
		t.Log.Info("can't wait end of program", "error", waitErr)
		return output.Result{Message: outputBuffer.Bytes(), Error: bencherr.ExecError(outputBuffer.String())}
	}

	if cmd.ProcessState.ExitCode() != 0 || !cmd.ProcessState.Success() {
		return output.Result{Message: outputBuffer.Bytes(), Error: bencherr.ExecError(outputBuffer.String())}
	}
	if releaseErr := cmd.Process.Release(); releaseErr != nil {
		t.Log.Info("can't release job process", "error", releaseErr)
	}
	return output.Result{Message: outputBuffer.Bytes(), Error: nil}
}
