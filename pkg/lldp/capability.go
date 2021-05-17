package lldp

import (
	"strconv"
	"strings"

	"github.com/onmetal/inventory/pkg/utils"
)

type Capability string

const (
	CLLDPOtherCapability             = "Other"
	CLLDPRepeaterCapability          = "Repeater"
	CLLDPBridgeCapability            = "Bridge"
	CLLDPWLANAccessPointCapability   = "WLAN Access Point"
	CLLDPRouterCapability            = "Router"
	CLLDPTelephoneCapability         = "Telephone"
	CLLDPDOCSISCableDeviceCapability = "DOCSIS cable device"
	CLLDPStationCapability           = "Station"
	CLLDPCustomerVLANCapability      = "Customer VLAN"
	CLLDPServiceVLANCapability       = "Service VLAN"
	CLLDPTwoPortMACRelayCapability   = "Two-port MAC Relay (TPMR)"
)

var CCapabilities = []Capability{
	CLLDPOtherCapability,
	CLLDPRepeaterCapability,
	CLLDPBridgeCapability,
	CLLDPWLANAccessPointCapability,
	CLLDPRouterCapability,
	CLLDPTelephoneCapability,
	CLLDPDOCSISCableDeviceCapability,
	CLLDPStationCapability,
	CLLDPCustomerVLANCapability,
	CLLDPServiceVLANCapability,
	CLLDPTwoPortMACRelayCapability,
}

func GetCapabilities(caps string) ([]Capability, error) {
	capabilities := make([]Capability, 0)
	for _, i := range strings.Split(caps, " ") {
		if i == "00" {
			continue
		}
		if parsed, err := strconv.ParseUint(i, 16, 8); err == nil {
			bitsList := make([]int, 0)
			utils.GetBitsList(uint8(parsed), &bitsList)
			for _, v := range bitsList {
				capabilities = append(capabilities, CCapabilities[v])
			}
		} else {
			return nil, err
		}
	}
	return capabilities, nil
}
