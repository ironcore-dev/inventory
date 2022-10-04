package nic

import (
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
