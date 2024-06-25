// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
