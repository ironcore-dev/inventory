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
