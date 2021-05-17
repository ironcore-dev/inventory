package lldp

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
