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

package frame

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"time"

	"github.com/mdlayher/lldp"
	"github.com/pkg/errors"
)

type Frame struct {
	InterfaceID         string
	ChassisID           string
	SystemName          string
	SystemDescription   string
	Capabilities        []Capability
	EnabledCapabilities []Capability
	PortID              string
	PortDescription     string
	ManagementAddresses []string
	TTL                 time.Duration
}

func (f *Frame) setChassisID(chassisID *lldp.ChassisID) error {
	switch chassisID.Subtype {
	case lldp.ChassisIDSubtypeChassisComponenent:
		fallthrough
	case lldp.ChassisIDSubtypeInterfaceAlias:
		fallthrough
	case lldp.ChassisIDSubtypePortComponent:
		fallthrough
	case lldp.ChassisIDSubtypeInterfaceName:
		fallthrough
	case lldp.ChassisIDSubtypeLocallyAssigned:
		f.ChassisID = string(chassisID.ID)
	case lldp.ChassisIDSubtypeMACAddress:
		chassisId, err := idBytesToMac(chassisID.ID)
		if err != nil {
			return errors.Wrap(err, "unable to convert chassis ID to MAC address")
		}
		f.ChassisID = chassisId
	case lldp.ChassisIDSubtypeNetworkAddress:
		chassisId, err := idBytesToNetworkAddress(chassisID.ID)
		if err != nil {
			return errors.Wrap(err, "unable to convert chassis ID to MAC address")
		}
		f.ChassisID = chassisId
	default:
		// fallback
		f.ChassisID = hex.EncodeToString(chassisID.ID)
	}

	return nil
}

func (f *Frame) setPortID(portID *lldp.PortID) error {
	switch portID.Subtype {
	case lldp.PortIDSubtypeInterfaceAlias:
		fallthrough
	case lldp.PortIDSubtypePortComponent:
		fallthrough
	case lldp.PortIDSubtypeInterfaceName:
		fallthrough
	case lldp.PortIDSubtypeLocallyAssigned:
		f.PortID = string(portID.ID)
	case lldp.PortIDSubtypeMACAddress:
		portId, err := idBytesToMac(portID.ID)
		if err != nil {
			return errors.Wrap(err, "unable to convert port ID to MAC address")
		}
		f.PortID = portId
	case lldp.PortIDSubtypeNetworkAddress:
		portId, err := idBytesToNetworkAddress(portID.ID)
		if err != nil {
			return errors.Wrap(err, "unable to convert port ID to network address")
		}
		f.PortID = portId
	default:
		// fallback
		f.PortID = hex.EncodeToString(portID.ID)
	}

	return nil
}

func (f *Frame) setOptional(tlv *lldp.TLV) error {
	switch tlv.Type {
	case lldp.TLVTypeSystemName:
		f.setSystemName(tlv.Value)
	case lldp.TLVTypeSystemDescription:
		f.setSystemDescription(tlv.Value)
	case lldp.TLVTypePortDescription:
		f.setPortDescription(tlv.Value)
	case lldp.TLVTypeSystemCapabilities:
		f.setSystemCapability(tlv.Value)
	case lldp.TLVTypeManagementAddress:
		if err := f.setManagementAddress(tlv.Value); err != nil {
			return errors.Wrapf(err, "unable to set management address")
		}
		// there is also lldp.TLVTypeOrganizationSpecific
		// not collecting it for now
	default:
		return errors.Errorf("unhandled TLV type %d", tlv.Type)
	}

	return nil
}

func (f *Frame) setSystemName(val []byte) {
	f.SystemName = string(val)
}

func (f *Frame) setSystemDescription(val []byte) {
	f.SystemDescription = string(val)
}

func (f *Frame) setPortDescription(val []byte) {
	f.PortDescription = string(val)
}

func (f *Frame) setSystemCapability(val []byte) {
	capabilities := binary.BigEndian.Uint16(val[0:2])
	enabledCapabilities := binary.BigEndian.Uint16(val[2:4])

	for i, capability := range CCapabilities {
		idx := uint16(1 << i)
		capable := capabilities & idx
		if capable != 0 {
			f.addCapability(capability)
		}

		enabled := enabledCapabilities & idx
		if enabled != 0 {
			f.addEnabledCapability(capability)
		}
	}
}

func (f *Frame) addCapability(capability Capability) {
	if f.Capabilities == nil {
		f.Capabilities = make([]Capability, 0)
	}

	f.Capabilities = append(f.Capabilities, capability)
}

func (f *Frame) addEnabledCapability(capability Capability) {
	if f.EnabledCapabilities == nil {
		f.EnabledCapabilities = make([]Capability, 0)
	}

	f.EnabledCapabilities = append(f.EnabledCapabilities, capability)
}

func (f *Frame) setManagementAddress(val []byte) error {
	var addressBytes []byte

	valLen := len(val)
	if valLen < 6 {
		return errors.Errorf("expected to have at least 6 bytes, but got %d", valLen)
	}

	switch val[1] {
	case 1:
		addressBytes = val[2:6]
	case 2:
		addressBytes = val[2:18]
	default:
		return errors.Errorf("unhandled address type %d", val[1])
	}

	ip := net.IP(addressBytes)

	f.addManagementAddress(ip.String())

	return nil
}

func (f *Frame) addManagementAddress(address string) {
	if f.ManagementAddresses == nil {
		f.ManagementAddresses = make([]string, 0)
	}

	f.ManagementAddresses = append(f.ManagementAddresses, address)
}
