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
	"github.com/pkg/errors"

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

func (s *DeviceSvc) GetDevice(thePath string, name string) (*Device, error) {
	nic := &Device{
		Name: name,
	}

	defs := []func(string) error{
		nic.defPCIAddress,
		nic.defAddressAssignType,
		nic.defAddress,
		nic.defAddressLength,
		nic.defBroadcast,
		nic.defCarrier,
		nic.defCarrierChanges,
		nic.defCarrierDownCount,
		nic.defCarrierUpCount,
		nic.defDevID,
		nic.defDevPort,
		nic.defDormant,
		nic.defDuplex,
		nic.defFlags,
		nic.defInterfaceAlias,
		nic.defInterfaceIndex,
		nic.defInterfaceLink,
		nic.defLinkMode,
		nic.defMTU,
		nic.defNameAssignType,
		nic.defNetDevGroup,
		nic.defOperationalState,
		nic.defPhysicalPortID,
		nic.defPhysicalPortName,
		nic.defPhysicalSwitchID,
		nic.defSpeed,
		nic.defTesting,
		nic.defTransmitQueueLength,
		nic.defType,
	}

	for _, def := range defs {
		err := def(thePath)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to set Device property"))
		}
	}

	return nic, nil
}
