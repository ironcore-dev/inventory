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
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/file"
)

const (
	CNICDevicePCIAddressPath = "/device"

	CNICDeviceAddressAddressAssignTypePath = "/addr_assign_type"
	CNICDeviceAddressPath                  = "/address"
	CNICDeviceAddressLengthPath            = "/addr_len"
	CNICDeviceBroadcastPath                = "/broadcast"
	CNICDeviceCarrierPath                  = "/carrier"
	CNICDeviceCarrierChangesPath           = "/carrier_changes"
	CNICDeviceCarrierDownCountPath         = "/carrier_down_count"
	CNICDeviceCarrierUpCountPath           = "/carrier_up_count"
	CNICDeviceDevIDPath                    = "/dev_id"
	CNICDeviceDevPortPath                  = "/dev_port"
	CNICDeviceDormantPath                  = "/dormant"
	CNICDeviceDuplexPath                   = "/duplex"
	CNICDeviceFlagsPath                    = "/flags"
	CNICDeviceInterfaceAliasPath           = "/ifalias"
	CNICDeviceInterfaceIndexPath           = "/ifindex"
	CNICDeviceInterfaceLinkPath            = "/iflink"
	CNICDeviceLinkModePath                 = "/link_mode"
	CNICDeviceMTUPath                      = "/mtu"
	CNICDeviceNameAssignTypePath           = "/name_assign_type"
	CNICDeviceNetDevGroupPath              = "/netdev_group"
	CNICDeviceOperationalStatePath         = "/operstate"
	CNICDevicePhysicalPortIDPath           = "/phys_port_id"
	CNICDevicePhysicalPortNamePath         = "/phys_port_name"
	CNICDevicePhysicalSwitchIDPath         = "/phys_switch_id"
	CNICDeviceSpeedPath                    = "/speed"
	CNICDeviceTestingPath                  = "/testing"
	CNICDeviceTransmitQueueLengthPath      = "/tx_queue_len"
	CNICDeviceTypePath                     = "/type"

	CPermanentAddressAssignType               = "permanent address"
	CRandomlyGeneratedAddressAssignType       = "randomly generated"
	CStolenFromAnotherDeviceAddressAssignType = "stolen from another device"
	CSetUsingDevAddressAssignType             = "set using dev_set_mac_address"

	CDefaultLinkMode = "default link mode"
	CDormantLinkMode = "dormant link mode"

	CUnpredictableKernelNameAssignType = "enumerated by the kernel, possibly in an unpredictable way"
	CPredictableKernelNameAssignType   = "predictably named by the kernel"
	CUserspaceNameAssignType           = "named by userspace"
	CRenamedNameAssignType             = "renamed"

	CPCIAddressPattern = "^([0-9a-fA-F]{4}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}.[0-9]{1})|([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$"
)

var CPCIAddressRegex = regexp.MustCompile(CPCIAddressPattern)

type AddressAssignType string

type LinkMode string

type NameAssignType string

var CAddressAssignTypes = []AddressAssignType{
	CPermanentAddressAssignType,
	CRandomlyGeneratedAddressAssignType,
	CStolenFromAnotherDeviceAddressAssignType,
	CSetUsingDevAddressAssignType,
}

var CLinkModes = []LinkMode{
	CDefaultLinkMode,
	CDormantLinkMode,
}

var CNameAssignTypes = []NameAssignType{
	CUnpredictableKernelNameAssignType,
	CPredictableKernelNameAssignType,
	CUserspaceNameAssignType,
	CRenamedNameAssignType,
}

type Device struct {
	Name       string
	PCIAddress string

	AddressAssignType   AddressAssignType
	Address             string
	AddressLength       uint8
	Broadcast           string
	Carrier             bool
	CarrierChanges      uint32
	CarrierDownCount    uint32
	CarrierUpCount      uint32
	DevID               string
	DevPort             uint8
	Dormant             bool
	Duplex              string
	Flags               *Flags
	InterfaceAlias      string
	InterfaceIndex      uint32
	InterfaceLink       uint32
	LinkMode            LinkMode
	MTU                 uint16
	NameAssignType      NameAssignType
	NetDevGroup         int
	OperationalState    string
	PhysicalPortID      string
	PhysicalPortName    string
	PhysicalSwitchID    string
	Speed               uint32
	Testing             bool
	TransmitQueueLength uint32
	Type                Type
	Lanes               uint8
	FEC                 string
}

func (n *Device) defPCIAddress(thePath string) error {
	filePath := path.Join(thePath, CNICDevicePCIAddressPath)

	linkPath, err := filepath.EvalSymlinks(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable resolve symlink with path %s", filePath)
	}
	fileInfo, err := os.Stat(linkPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get stat for path %s", filePath)
	}

	// Device may use other bus, e.g. USB,
	// therefore we need to validate whether it is a PCI address or not first
	pciAddress := fileInfo.Name()
	isPCIAddress := CPCIAddressRegex.MatchString(pciAddress)
	if !isPCIAddress {
		pciAddress = ""
	}
	n.PCIAddress = pciAddress

	return nil
}

func (n *Device) defAddressAssignType(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceAddressAddressAssignTypePath)
	fileVal, err := file.ToInt(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get address assign type from %s", filePath)
	}

	if fileVal < 0 || fileVal >= len(CAddressAssignTypes) {
		return errors.Errorf("unexpected value %d for address assign type", fileVal)
	}
	n.AddressAssignType = CAddressAssignTypes[fileVal]

	return nil
}

func (n *Device) defAddress(thePath string) error {
	addressPath := path.Join(thePath, CNICDeviceAddressPath)
	addressString, err := file.ToString(addressPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get address string from %s", addressPath)
	}

	n.Address = addressString

	return nil
}

func (n *Device) defAddressLength(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceAddressLengthPath)
	fileVal, err := file.ToUint8(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get address length from %s", filePath)
	}

	n.AddressLength = fileVal

	return nil
}

func (n *Device) defBroadcast(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceBroadcastPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get broadcast address from %s", filePath)
	}

	n.Broadcast = fileVal

	return nil
}

func (n *Device) defCarrier(thePath string) error {
	carrierPath := path.Join(thePath, CNICDeviceCarrierPath)
	carrier, err := file.ToBool(carrierPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get carrier status from %s", carrierPath)
	}

	n.Carrier = carrier

	return nil
}

func (n *Device) defCarrierChanges(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceCarrierChangesPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get carrier changes address from %s", filePath)
	}

	n.CarrierChanges = fileVal

	return nil
}

func (n *Device) defCarrierDownCount(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceCarrierDownCountPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get carrier down count from %s", filePath)
	}

	n.CarrierChanges = fileVal

	return nil
}

func (n *Device) defCarrierUpCount(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceCarrierUpCountPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get carrier up count from %s", filePath)
	}

	n.CarrierChanges = fileVal

	return nil
}

func (n *Device) defDevID(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceDevIDPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get dev ID from %s", filePath)
	}

	n.DevID = fileVal

	return nil
}

func (n *Device) defDevPort(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceDevPortPath)
	fileVal, err := file.ToUint8(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get dev port from %s", filePath)
	}

	n.DevPort = fileVal

	return nil
}

func (n *Device) defDormant(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceDormantPath)
	fileVal, err := file.ToBool(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get dormant value from %s", filePath)
	}

	n.Dormant = fileVal

	return nil
}

func (n *Device) defDuplex(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceDuplexPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get duplex value from %s", filePath)
	}

	n.Duplex = fileVal

	return nil
}

func (n *Device) defFlags(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceFlagsPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get flags value from %s", filePath)
	}

	val, err := strconv.ParseUint(fileVal[2:], 16, 32)
	if err != nil {
		return errors.Wrapf(err, "unable to parse uint from from %s", fileVal)
	}

	n.Flags = NewFlags(uint32(val))

	return nil
}

func (n *Device) defInterfaceAlias(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceInterfaceAliasPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get interface alias value from %s", filePath)
	}

	n.InterfaceAlias = fileVal

	return nil
}

func (n *Device) defInterfaceIndex(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceInterfaceIndexPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get interface index value from %s", filePath)
	}

	n.InterfaceIndex = fileVal

	return nil
}

func (n *Device) defInterfaceLink(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceInterfaceLinkPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get interface link value from %s", filePath)
	}

	n.InterfaceLink = fileVal

	return nil
}

func (n *Device) defLinkMode(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceLinkModePath)
	fileVal, err := file.ToInt(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get link mode from %s", filePath)
	}

	if fileVal < 0 || fileVal >= len(CLinkModes) {
		return errors.Errorf("unexpected value %d for link mode", fileVal)
	}
	n.LinkMode = CLinkModes[fileVal]

	return nil
}

func (n *Device) defMTU(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceMTUPath)
	fileVal, err := file.ToUint16(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get MTU value from %s", filePath)
	}

	n.MTU = fileVal

	return nil
}

func (n *Device) defNameAssignType(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceNameAssignTypePath)
	fileVal, err := file.ToInt(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get name assign type from %s", filePath)
	}

	fileVal = fileVal - 1

	if fileVal < 0 || fileVal >= len(CNameAssignTypes) {
		return errors.Errorf("unexpected value %d for name assign type", fileVal)
	}
	n.NameAssignType = CNameAssignTypes[fileVal]

	return nil
}

func (n *Device) defNetDevGroup(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceNetDevGroupPath)
	fileVal, err := file.ToInt(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get network device group value from %s", filePath)
	}

	n.NetDevGroup = fileVal

	return nil
}

func (n *Device) defOperationalState(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceOperationalStatePath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get operational state value from %s", filePath)
	}

	n.OperationalState = fileVal

	return nil
}

func (n *Device) defPhysicalPortID(thePath string) error {
	filePath := path.Join(thePath, CNICDevicePhysicalPortIDPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get physical port ID value from %s", filePath)
	}

	n.PhysicalPortID = fileVal

	return nil
}

func (n *Device) defPhysicalPortName(thePath string) error {
	filePath := path.Join(thePath, CNICDevicePhysicalPortNamePath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get physical port name value from %s", filePath)
	}

	n.PhysicalPortName = fileVal

	return nil
}

func (n *Device) defPhysicalSwitchID(thePath string) error {
	filePath := path.Join(thePath, CNICDevicePhysicalSwitchIDPath)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get physical switch ID value from %s", filePath)
	}

	n.PhysicalSwitchID = fileVal

	return nil
}

func (n *Device) defSpeed(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceSpeedPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get speed value from %s", filePath)
	}

	n.Speed = fileVal

	return nil
}

func (n *Device) defTesting(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceTestingPath)
	fileVal, err := file.ToBool(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get testing value from %s", filePath)
	}

	n.Testing = fileVal

	return nil
}

func (n *Device) defTransmitQueueLength(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceTransmitQueueLengthPath)
	fileVal, err := file.ToUint32(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get transmit queue length value from %s", filePath)
	}

	n.TransmitQueueLength = fileVal

	return nil
}

func (n *Device) defType(thePath string) error {
	filePath := path.Join(thePath, CNICDeviceTypePath)
	fileVal, err := file.ToUint16(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get type value from %s", filePath)
	}

	theType, ok := CTypes[fileVal]
	if ok {
		n.Type = theType
	} else {
		n.Type = CTypes[0xffff]
	}

	return nil
}
