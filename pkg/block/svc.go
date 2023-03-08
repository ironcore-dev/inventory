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

package block

import (
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CSysBlockBasePath = "/sys/block"
)

type Svc struct {
	printer      *printer.Svc
	devSvc       *DeviceSvc
	sysBlockPath string
}

func NewSvc(printer *printer.Svc, devSvc *DeviceSvc, basePath string) *Svc {
	return &Svc{
		printer:      printer,
		devSvc:       devSvc,
		sysBlockPath: path.Join(basePath, CSysBlockBasePath),
	}
}

func (s *Svc) GetData() ([]Device, error) {
	blocks := make([]Device, 0)
	fileInfos, err := os.ReadDir(s.sysBlockPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of block devices")
	}

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		thePath := path.Join(s.sysBlockPath, name)
		block, err := s.devSvc.GetDevice(thePath, name)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to collect block device data for %s", thePath))
		}

		blocks = append(blocks, *block)
	}

	return blocks, nil
}
