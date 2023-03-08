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

package ipmi

import (
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CDevPath        = "/dev"
	CIPMIDevPattern = "ipmi\\d+"
)

var CIPMIDevRegexp = regexp.MustCompile(CIPMIDevPattern)

type Svc struct {
	printer     *printer.Svc
	ipmiInfoSvc *DeviceSvc
	devPath     string
}

func NewSvc(printer *printer.Svc, ipmiDevInfoSvc *DeviceSvc, basePath string) *Svc {
	return &Svc{
		printer:     printer,
		ipmiInfoSvc: ipmiDevInfoSvc,
		devPath:     path.Join(basePath, CDevPath),
	}
}

func (s *Svc) GetData() ([]Device, error) {
	devFolderContents, err := os.ReadDir(s.devPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read contents of %s", s.devPath)
	}

	infos := make([]Device, 0)
	for _, dev := range devFolderContents {
		devName := dev.Name()

		matches := CIPMIDevRegexp.MatchString(devName)

		if !matches {
			continue
		}

		thePath := path.Join(s.devPath, devName)
		info, err := s.ipmiInfoSvc.GetDevice(thePath)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unabale to obtain IPMI device info"))
			continue
		}

		infos = append(infos, *info)
	}

	return infos, nil
}
