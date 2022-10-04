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
