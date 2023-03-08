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
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/u-root/u-root/pkg/ipmi"
)

const (
	CIPMIIOCtlChannel = 1

	CIPMIIOCtlSetInProgressFlag   = 0
	CIPMIIOCtlIPAddressFlag       = 3
	CIPMIIOCtlIPAddressSourceFlag = 4
	CIPMIIOCtlMACAddressFlag      = 5

	CIPMISetInProgressSetCompleteStatus   = "Set Complete"
	CIPMISetInProgressSetInProgressStatus = "Set In Progress"
	CIPMISetInProgressCommitWriteStatus   = "Commit Write"
	CIPMISetInProgressReservedStatus      = "Reserved"
	CIPMISetInProgressUnknownStatus       = "Unknown"

	CIPMIIPAddressSourceUnspecified = "Unspecified"
	CIPMIIPAddressSourceStatic      = "Static Address"
	CIPMIIPAddressSourceDHCP        = "DHCP Address"
	CIPMIIPAddressSourceBIOS        = "BIOS Assigned Address"
	CIPMIIPAddressSourceOther       = "Other"

	CIPMIAdditionalSensorDeviceSupport        = "Sensor Device"
	CIPMIAdditionalSDRRepositoryDeviceSupport = "SDR Repository Device"
	CIPMIAdditionalSELDeviceSupport           = "SEL Device"
	CIPMIAdditionalFRUInventoryDeviceSupport  = "FRU Inventory Device"
	CIPMIAdditionalIPMBEventReceiverSupport   = "IPMB Event Receiver"
	CIPMIAdditionalIPMBEventGeneratorSupport  = "IPMB Event Generator"
	CIPMIAdditionalBridgeSupport              = "Bridge"
	CIPMIAdditionalChassisDeviceSupport       = "Chassis Device"
)

type SetInProgressStatus string

var CSetInProgressStatuses = []SetInProgressStatus{
	CIPMISetInProgressSetCompleteStatus,
	CIPMISetInProgressSetInProgressStatus,
	CIPMISetInProgressCommitWriteStatus,
	CIPMISetInProgressReservedStatus,
}

type IPAddressSource string

var CIPAddressSources = []IPAddressSource{
	CIPMIIPAddressSourceUnspecified,
	CIPMIIPAddressSourceStatic,
	CIPMIIPAddressSourceDHCP,
	CIPMIIPAddressSourceBIOS,
}

type AdditionalDeviceSupport string

var CAdditionalDeviceSupportList = []AdditionalDeviceSupport{
	CIPMIAdditionalSensorDeviceSupport,
	CIPMIAdditionalSDRRepositoryDeviceSupport,
	CIPMIAdditionalSELDeviceSupport,
	CIPMIAdditionalFRUInventoryDeviceSupport,
	CIPMIAdditionalIPMBEventReceiverSupport,
	CIPMIAdditionalIPMBEventGeneratorSupport,
	CIPMIAdditionalBridgeSupport,
	CIPMIAdditionalChassisDeviceSupport,
}

type Device struct {
	ID                      uint8
	Revision                uint8
	FirmwareRevision        string
	IPMIVersion             string
	ManufacturerID          string
	ProductID               string
	DeviceAvailable         bool
	ProvidesDeviceSDRs      bool
	AdditionalDeviceSupport []AdditionalDeviceSupport
	AuxFirmwareRevInfo      []string

	SetInProgress   SetInProgressStatus
	IPAddressSource IPAddressSource
	IPAddress       string
	MACAddress      string
}

func (i *Device) defDevice(conn *ipmi.IPMI) error {
	deviceInfo, err := conn.GetDeviceID()
	if err != nil {
		return err
	}

	i.ID = deviceInfo.DeviceID
	i.Revision = deviceInfo.DeviceRevision & 0x0f
	i.FirmwareRevision = fmt.Sprintf("%d.%02x", deviceInfo.FwRev1&0x3f, deviceInfo.FwRev2)
	i.IPMIVersion = fmt.Sprintf("%x.%x", deviceInfo.IpmiVersion&0x0f, (deviceInfo.IpmiVersion&0x0f)>>4)

	if deviceInfo.AdtlDeviceSupport != 0 {
		i.AdditionalDeviceSupport = make([]AdditionalDeviceSupport, 0)
	}

	manufacturerID := uint32(deviceInfo.ManufacturerID[2]) << 16
	manufacturerID |= uint32(deviceInfo.ManufacturerID[1]) << 8
	manufacturerID |= uint32(deviceInfo.ManufacturerID[0])
	i.ManufacturerID = fmt.Sprintf("%d (0x%04X)", manufacturerID, manufacturerID)

	productID := uint16(deviceInfo.ProductID[1]) << 8
	productID |= uint16(deviceInfo.ProductID[1])
	i.ProductID = fmt.Sprintf("%d (0x%04X)", productID, productID)

	i.DeviceAvailable = (^deviceInfo.FwRev1 & 0x80) != 0
	i.ProvidesDeviceSDRs = deviceInfo.DeviceRevision&0x80 != 0

	for shift, sup := range CAdditionalDeviceSupportList {
		idx := byte(1 << shift)
		val := deviceInfo.AdtlDeviceSupport & idx

		if val != 0 {
			i.AdditionalDeviceSupport = append(i.AdditionalDeviceSupport, sup)
		}
	}

	i.AuxFirmwareRevInfo = make([]string, len(deviceInfo.AuxFwRev))

	for idx, val := range deviceInfo.AuxFwRev {
		i.AuxFirmwareRevInfo[idx] = fmt.Sprintf("0x%02x", val)
	}

	return nil
}

func (i *Device) defSetInProgress(conn *ipmi.IPMI) error {
	bytes, err := conn.GetLanConfig(CIPMIIOCtlChannel, CIPMIIOCtlSetInProgressFlag)
	if err != nil {
		return errors.Wrap(err, "unable to get set in progress")
	}

	if len(bytes) < 3 {
		return errors.Wrap(err, "unable to get set in progress")
	}

	valIdx := int(bytes[2])

	if valIdx >= len(CSetInProgressStatuses) {
		i.SetInProgress = CIPMISetInProgressUnknownStatus
	} else {
		i.SetInProgress = CSetInProgressStatuses[valIdx]
	}

	return nil
}

func (i *Device) defIPAddressSource(conn *ipmi.IPMI) error {
	bytes, err := conn.GetLanConfig(CIPMIIOCtlChannel, CIPMIIOCtlIPAddressSourceFlag)
	if err != nil {
		return errors.Wrap(err, "unable to get IP address source")
	}

	if len(bytes) < 3 {
		return errors.New("unable to obtain IP address source")
	}

	valIdx := int(bytes[2])

	if valIdx >= len(CIPAddressSources) {
		i.IPAddressSource = CIPMIIPAddressSourceOther
	} else {
		i.IPAddressSource = CIPAddressSources[valIdx]
	}

	return nil
}

func (i *Device) defIPAddress(conn *ipmi.IPMI) error {
	bytes, err := conn.GetLanConfig(CIPMIIOCtlChannel, CIPMIIOCtlIPAddressFlag)
	if err != nil {
		return errors.Wrap(err, "unable to obtain IP address")
	}

	if len(bytes) < 6 {
		return errors.New("unable to obtain IP address")
	}

	ipBytes := bytes[2:6]
	ip := net.IP(ipBytes)
	i.IPAddress = ip.String()

	return nil
}

func (i *Device) defMACAddress(conn *ipmi.IPMI) error {
	bytes, err := conn.GetLanConfig(CIPMIIOCtlChannel, CIPMIIOCtlMACAddressFlag)
	if err != nil {
		return errors.Wrap(err, "unable to obtain MAC address")
	}

	if len(bytes) < 8 {
		return errors.New("unable to obtain MAC address")
	}

	macBytes := bytes[2:8]
	mac := net.HardwareAddr(macBytes)
	i.MACAddress = mac.String()

	return nil
}
