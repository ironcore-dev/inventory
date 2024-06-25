// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package ipmi

import (
	"os"

	"github.com/pkg/errors"
	"github.com/u-root/u-root/pkg/ipmi"

	"github.com/onmetal/inventory/pkg/printer"
)

type DeviceSvc struct {
	printer *printer.Svc
}

func NewDeviceSvc(printer *printer.Svc) *DeviceSvc {
	return &DeviceSvc{
		printer: printer,
	}
}

func (s *DeviceSvc) GetDevice(thePath string) (*Device, error) {
	f, err := os.OpenFile(thePath, os.O_RDWR, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open IPMI device file %s", thePath)
	}

	conn := &ipmi.IPMI{
		File: f,
	}
	defer func() {
		if err := conn.Close(); err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to close file %s", thePath))
		}
	}()

	info := &Device{}

	defs := []func(*ipmi.IPMI) error{
		info.defDevice,
		info.defSetInProgress,
		info.defIPAddressSource,
		info.defIPAddress,
		info.defMACAddress,
	}

	for _, def := range defs {
		err := def(conn)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to set IPMI device info property"))
		}
	}

	return info, nil
}
