// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package block

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CStatPath = "/stat"
)

type DeviceStatSvc struct {
	printer *printer.Svc
}

func NewDeviceStatSvc(printer *printer.Svc) *DeviceStatSvc {
	return &DeviceStatSvc{
		printer: printer,
	}
}

func (s *DeviceStatSvc) GetDeviceStat(thePath string) (*DeviceStat, error) {
	statPath := path.Join(thePath, CStatPath)
	contents, err := os.ReadFile(statPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file %s", statPath)
	}

	stringContents := string(contents)
	trimmedStringContents := strings.TrimSpace(stringContents)

	fields := strings.Fields(trimmedStringContents)

	statVals := make([]uint64, len(fields))
	for i, field := range fields {
		val, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to convert to uint64 %s", field))
			val = 0
		}

		statVals[i] = val
	}

	stat := &DeviceStat{}

	// linux kernel doc states that there are 11 fields
	// and underneath there is a table for 17
	// guess, we need to check this for the backward compatibility
	for i, val := range statVals {
		if err := stat.setByIndex(i, val); err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to set value %d on index %d", val, i))
		}
	}

	return stat, nil
}
