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
	CPCIDeviceAddressPattern = "[[:xdigit:]]{4}:[[:xdigit:]]{2}:[[:xdigit:]]{2}.[[:xdigit:]]"
)

var CPCIDeviceAddressRegexp = regexp.MustCompile(CPCIDeviceAddressPattern)

type BusSvc struct {
	printer      *printer.Svc
	pciDeviceSvc *DeviceSvc
}

func NewBusSvc(printer *printer.Svc, pciDevSvc *DeviceSvc) *BusSvc {
	return &BusSvc{
		printer:      printer,
		pciDeviceSvc: pciDevSvc,
	}
}

func (s *BusSvc) GetBus(thePath string, id string) (*Bus, error) {
	pciDevFolders, err := os.ReadDir(thePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of block devices")
	}

	devices := make([]Device, 0)
	for _, pciDevFolder := range pciDevFolders {
		name := pciDevFolder.Name()

		if !pciDevFolder.IsDir() {
			continue
		}

		match := CPCIDeviceAddressRegexp.MatchString(name)
		if !match {
			continue
		}

		pciDevPath := path.Join(thePath, name)

		device, err := s.pciDeviceSvc.GetDevice(pciDevPath, name)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to collect info for device %s", name))
			continue
		}

		devices = append(devices, *device)
	}

	return NewBus(id, devices), nil
}
