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

package nic

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	switchConstants "github.com/ironcore-dev/metal/pkg/constants"

	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/redis"
	"github.com/onmetal/inventory/pkg/utils"
)

const (
	CNICDevicePath = "/sys/class/net"
)

type Svc struct {
	printer    *printer.Svc
	nicDevSvc  *DeviceSvc
	nicDevPath string
	hostSvc    *host.Svc
	redisSvc   *redis.Svc
}

func NewSvc(printer *printer.Svc, nicDevSvc *DeviceSvc, hostSvc *host.Svc, redisSvc *redis.Svc, basePath string) *Svc {
	return &Svc{
		printer:    printer,
		nicDevSvc:  nicDevSvc,
		hostSvc:    hostSvc,
		redisSvc:   redisSvc,
		nicDevPath: path.Join(basePath, CNICDevicePath),
	}
}

func (s *Svc) GetData() ([]Device, error) {
	hostInfo, err := s.hostSvc.GetData()
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "failed to collect host info"))
	}

	nicFolders, err := os.ReadDir(s.nicDevPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of nic folders")
	}

	var nics []Device
	for _, nicFolder := range nicFolders {
		fName := nicFolder.Name()
		thePath := path.Join(s.nicDevPath, fName)
		nic, err := s.nicDevSvc.GetDevice(thePath, fName)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to collect Device data"))
			continue
		}
		if hostInfo.Type == utils.CSwitchType && strings.HasPrefix(fName, "Ethernet") {
			info, err := s.redisSvc.GetPortAdditionalInfo(fName)
			if err != nil {
				s.printer.VErr(errors.Wrap(err, "unable to collect additional Device data from Redis"))
				continue
			}
			nic.Lanes = uint8(len(strings.Split(info[redis.CPortLanes], ",")))
			nic.FEC = info[redis.CPortFec]
			if nic.FEC == "" {
				nic.FEC = switchConstants.FECNone
			}
			speed, err := strconv.Atoi(info[redis.CPortSpeed])
			if err != nil {
				s.printer.VErr(errors.Wrap(err, "unable to collect additional Device data from Redis"))
				continue
			}
			nic.Speed = uint32(speed)
		}
		nics = append(nics, *nic)
	}

	return nics, nil
}
