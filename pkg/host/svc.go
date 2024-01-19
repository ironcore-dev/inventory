// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package host

import (
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/utils"
)

type Info struct {
	Type string
	Name string
}

type Svc struct {
	printer           *printer.Svc
	switchVersionPath string
}

func NewSvc(printer *printer.Svc, basePath string) *Svc {
	return &Svc{
		printer:           printer,
		switchVersionPath: path.Join(basePath, utils.CVersionFilePath),
	}
}

func (s *Svc) GetData() (*Info, error) {
	hostType, err := getHostType(s.switchVersionPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine host type")
	}

	info := Info{}
	name, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hostname")
	}
	info.Name = name
	info.Type = hostType
	return &info, nil
}

func getHostType(versionFile string) (string, error) {
	//todo: determining how to check host type without checking files
	if _, err := os.Stat(versionFile); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		} else {
			return utils.CMachineType, nil
		}
	}
	return utils.CSwitchType, nil
}
