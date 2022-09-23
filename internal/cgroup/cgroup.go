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

package cgroup

import (
	"os"
	"path/filepath"
	"runtime"

	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"

	bencherr "github.com/onmetal/inventory/internal/errors"
)

type Manager interface {
	Add(pid int) error
	Delete() error
}

const (
	defaultRootPath         = "/"
	defaultCgroupMountPoint = "/sys/fs/cgroup"
	v2ControllerPath        = "/sys/fs/cgroup/cgroup.controllers"
)

func NewWithLimits(name string, resources *conf.Resources) (Manager, error) {
	root := defaultRootPath
	if os.Getenv("ROOT") != "" {
		root = os.Getenv("ROOT")
	}
	mountPoint := root + defaultCgroupMountPoint
	controllerPath := filepath.Join(root, v2ControllerPath)
	switch operationSystem := runtime.GOOS; operationSystem {
	case "linux":
		if !isFileExist(controllerPath) {
			return newV1(name, mountPoint, resources)
		}
		return newV2(name, mountPoint, resources)
	default:
		return nil, bencherr.NotSupportedOS(operationSystem)
	}
}

func isFileExist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
