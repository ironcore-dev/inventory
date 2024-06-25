// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
