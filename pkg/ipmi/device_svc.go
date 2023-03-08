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
