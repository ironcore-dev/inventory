// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package lldp

import (
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/lldp/frame"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/redis"
	"github.com/onmetal/inventory/pkg/utils"
)

const (
	CLLDPPath = "/run/systemd/netif/lldp"
)

type Svc struct {
	printer      *printer.Svc
	frameInfoSvc *frame.Svc
	hostSvc      *host.Svc
	redisSvc     *redis.Svc
	lldpPath     string
}

func NewSvc(printer *printer.Svc, frameInfoSvc *frame.Svc, hostSvc *host.Svc, redisSvc *redis.Svc, basePath string) *Svc {
	return &Svc{
		printer:      printer,
		frameInfoSvc: frameInfoSvc,
		hostSvc:      hostSvc,
		redisSvc:     redisSvc,
		lldpPath:     path.Join(basePath, CLLDPPath),
	}
}

func (s *Svc) GetData() ([]frame.Frame, error) {
	frameInfos := make([]frame.Frame, 0)

	hostInfo, err := s.hostSvc.GetData()
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "failed to collect host info"))
	}

	switch hostInfo.Type {
	case utils.CMachineType:
		frameFiles, err := os.ReadDir(s.lldpPath)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get list of frame files")
		}
		// iterate over /run/systemd/netif/lldp/%i
		for _, frameFile := range frameFiles {
			fName := frameFile.Name()
			filePath := path.Join(s.lldpPath, fName)
			info, err := s.frameInfoSvc.GetFrame(fName, filePath)
			if err != nil {
				s.printer.VErr(errors.Errorf("unable to collect LLDP info for interface idx %s", fName))
				continue
			}
			frameInfos = append(frameInfos, *info)
		}
	case utils.CSwitchType:
		frames, err := s.redisSvc.GetFrames()
		if err != nil {
			return nil, errors.Wrap(err, "unable to process redis lldp data")
		}
		frameInfos = append(frameInfos, frames...)
	}
	return frameInfos, nil
}
