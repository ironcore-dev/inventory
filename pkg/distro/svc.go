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

package distro

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/utils"
)

type Distro struct {
	BuildVersion  string
	DebianVersion string
	KernelVersion string
	AsicType      string
	CommitID      string
	BuildDate     string
	BuildNumber   uint32
	BuildBy       string
}

type Svc struct {
	printer           *printer.Svc
	hostSvc           *host.Svc
	switchVersionPath string
}

func NewSvc(printer *printer.Svc, hostSvc *host.Svc, basePath string) *Svc {
	return &Svc{
		printer:           printer,
		hostSvc:           hostSvc,
		switchVersionPath: path.Join(basePath, utils.CVersionFilePath),
	}
}

func (s *Svc) GetData() (*Distro, error) {
	distro := Distro{}
	rawInfo := make(map[string]interface{})
	hostInfo, err := s.hostSvc.GetData()
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "failed to collect host info"))
	}
	switch hostInfo.Type {
	case utils.CSwitchType:
		sonicInfo, err := os.ReadFile(s.switchVersionPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read SONiC version file")
		}
		err = yaml.Unmarshal(sonicInfo, &rawInfo)
		if err != nil {
			return nil, errors.Wrap(err, "failed to collect SONiC version")
		}
		err = convertMapStruct(&distro, rawInfo)
		if err != nil {
			return nil, errors.Wrap(err, "failed to process SONiC version")
		}
		// todo: case utils.CMachineType:
	}
	return &distro, nil
}

func convertMapStruct(obj *Distro, m map[string]interface{}) error {
	for k, v := range m {
		m[strings.Replace(k, "_", "", 1)] = v
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}
