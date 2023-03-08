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
