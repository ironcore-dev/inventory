// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
