// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package pci

import (
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/file"
	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CPCIDeviceVendorPath          = "/vendor"
	CPCIDeviceTypePath            = "/device"
	CPCIDeviceSubsystemDevicePath = "/subsystem_device"
	CPCIDeviceSubsystemVendorPath = "/subsystem_vendor"

	CPCIDeviceClassPath = "/class"

	CPCIDeviceClassPattern = "0x([[:xdigit:]]{2})([[:xdigit:]]{2})([[:xdigit:]]{2})"
)

var CPCIDeviceClassRegexp = regexp.MustCompile(CPCIDeviceClassPattern)

type DeviceSvc struct {
	ids     *IDs
	printer *printer.Svc
}

func NewDeviceSvc(printer *printer.Svc, ids *IDs) *DeviceSvc {
	return &DeviceSvc{
		ids:     ids,
		printer: printer,
	}
}

func (s *DeviceSvc) GetDevice(basePath string, addr string) (*Device, error) {
	device := &Device{
		Address: addr,
	}

	err := s.setVendor(device, basePath)
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to resolve vendor branch"))
	}

	err = s.setClass(device, basePath)
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to resolve class branch"))
	}

	return device, nil
}

func (s *DeviceSvc) setVendor(dev *Device, thePath string) error {
	vendorPath := path.Join(thePath, CPCIDeviceVendorPath)
	vendorString, err := file.ToString(vendorPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get vendor string from %s", vendorPath)
	}

	vendorVal := vendorString[2:]
	vendor, ok := s.ids.Vendors[vendorVal]
	if !ok {
		return errors.Errorf("unknown vendor id %s", vendorVal)
	}

	dev.Vendor = &DeviceVendor{
		ID:   vendor.ID,
		Name: vendor.Name,
	}

	if len(vendor.Devices) == 0 {
		return nil
	}

	if err := s.setType(dev, thePath, &vendor); err != nil {
		return errors.Wrapf(err, "unable to resolve device/type branch for vendor %s, %s", vendor.ID, vendor.Name)
	}

	return nil
}

func (s *DeviceSvc) setType(dev *Device, thePath string, vendor *Vendor) error {
	typePath := path.Join(thePath, CPCIDeviceTypePath)
	typeString, err := file.ToString(typePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get device/type string from %s", typePath)
	}

	typeVal := typeString[2:]
	theType, ok := vendor.Devices[typeVal]
	if !ok {
		return errors.Errorf("unknown device/type id %s", typeVal)
	}

	dev.Type = &DeviceType{
		ID:   theType.ID,
		Name: theType.Name,
	}

	if len(theType.Subsystems) == 0 {
		return nil
	}

	if err := s.setSubsystem(dev, thePath, &theType); err != nil {
		return errors.Wrapf(err, "unable to resolve subsystem branch for device %s, %s", theType.ID, theType.Name)
	}

	return nil
}

func (s *DeviceSvc) setSubsystem(dev *Device, thePath string, theType *Type) error {
	subDevicePath := path.Join(thePath, CPCIDeviceSubsystemDevicePath)
	subDeviceString, err := file.ToString(subDevicePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get subsystem device/type string from %s", subDevicePath)
	}
	subDeviceVal := subDeviceString[2:]

	subVendorPath := path.Join(thePath, CPCIDeviceSubsystemVendorPath)
	subVendorString, err := file.ToString(subVendorPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get subsystem vendor string from %s", subVendorPath)
	}
	subVendorVal := subVendorString[2:]

	subsystem, ok := theType.Subsystems[subVendorVal+subDeviceVal]
	if !ok {
		return errors.Errorf("unknown subsystem id %s %s", subVendorVal, subDeviceVal)
	}

	dev.Subtype = &DeviceSubtype{
		ID:   subsystem.SubdeviceID,
		Name: subsystem.Name,
	}

	subvendor, ok := s.ids.Vendors[subVendorVal]
	if !ok {
		return errors.Errorf("unknown subsystem vendor id %s", subVendorVal)
	}

	dev.Subvendor = &DeviceVendor{
		ID:   subvendor.ID,
		Name: subvendor.Name,
	}

	return nil
}

func (s *DeviceSvc) setClass(dev *Device, thePath string) error {
	classPath := path.Join(thePath, CPCIDeviceClassPath)
	classString, err := file.ToString(classPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get vendor string from %s", classPath)
	}

	groups := CPCIDeviceClassRegexp.FindStringSubmatch(classString)
	classVal := groups[1]
	subclassVal := groups[2]
	ifaceVal := groups[3]

	class, ok := s.ids.Classes[classVal]
	if !ok {
		return errors.Errorf("unknown device class %s", classVal)
	}
	dev.Class = &DeviceClass{
		ID:   class.ID,
		Name: class.Name,
	}

	if len(class.Subclasses) == 0 {
		return nil
	}

	subclass, ok := class.Subclasses[subclassVal]
	if !ok {
		return errors.Errorf("unknown device subclass %s in class %s %s", subclassVal, class.ID, class.Name)
	}
	dev.Subclass = &DeviceSubclass{
		ID:   subclass.ID,
		Name: subclass.Name,
	}

	if len(subclass.ProgrammingInterfaces) == 0 {
		return nil
	}

	iface, ok := subclass.ProgrammingInterfaces[ifaceVal]
	if !ok {
		return errors.Errorf("unknown device programming interface %s in subclass %s %s in class %s %s",
			ifaceVal, subclass.ID, subclass.Name, class.ID, class.Name)
	}
	dev.ProgrammingInterface = &DeviceProgrammingInterface{
		ID:   iface.ID,
		Name: iface.Name,
	}

	return nil
}
