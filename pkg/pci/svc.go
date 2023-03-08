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

package pci

import (
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CDevicesPath = "/sys/devices"

	CPCIBusIDPattern = "pci(\\d{4}:\\d{2})"
)

var CPCIBusIDRegexp = regexp.MustCompile(CPCIBusIDPattern)

type Svc struct {
	printer     *printer.Svc
	pciBusSvc   *BusSvc
	devicesPath string
}

func NewSvc(printer *printer.Svc, pciBusSvc *BusSvc, basePath string) *Svc {
	return &Svc{
		printer:     printer,
		pciBusSvc:   pciBusSvc,
		devicesPath: path.Join(basePath, CDevicesPath),
	}
}

func (s *Svc) GetData() ([]Bus, error) {
	deviceFolders, err := os.ReadDir(s.devicesPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of device folders")
	}

	var buses []Bus
	for _, deviceFolder := range deviceFolders {
		fName := deviceFolder.Name()

		groups := CPCIBusIDRegexp.FindStringSubmatch(fName)

		if len(groups) != 2 {
			continue
		}

		pciBusPath := path.Join(s.devicesPath, fName)
		bus, err := s.pciBusSvc.GetBus(pciBusPath, groups[1])
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to collect PCI bus %s data", groups[1]))
			continue
		}

		buses = append(buses, *bus)
	}

	return buses, nil
}
