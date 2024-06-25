// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package nic

type Type string

const (
	CNICNETROMType               = "from KA9Q: NET/ROM pseudo"
	CNICEthernetType             = "Ethernet 10Mbps"
	CNICExperimentalEthernetType = "Experimental Ethernet"
	CNICAX25Type                 = "AX.25 Level 2"
	CNICPRONetType               = "PROnet token ring"
	CNICChaosNetType             = "Chaosnet"
	CNICIEEE802Type              = "IEEE 802.2 Ethernet/TR/TB"
	CNICARCNetType               = "ARCnet"
	CNICAppleTalkType            = "APPLEtalk"
	CNICDLCIType                 = "Frame Relay DLC"
	CNICATMType                  = "ATM"
	CNICMetricomType             = "Metricom STRIP (new IANA id)"
	CNICIEEE1394Type             = "IEEE 1394 IPv4 - RFC 2734"
	CNICEUI64Type                = "EUI-64"
	CNICInfiniBandType           = "InfiniBand"
	CNICSlipType                 = "Slip"
	CNICCSlipType                = "CSlip"
	CNICSlip6Type                = "Slip6"
	CNICCSlip6Type               = "CSlip6"
	CNICRSRVDType                = "Notional KISS type"
	CNICAdaptType                = "Adapt"
	CNICRoseType                 = "Rose"
	CNICX25Type                  = "CCITT X.25"
	CNICHWX25Type                = "Boards with X.25 in firmware"
	CNICCANType                  = "Controller Area Network"
	CNICPPPType                  = "PPP"
	CNICCiscoHDLCType            = "Cisco HDLC"
	CNICLAPBType                 = "LAPB"
	CNICDDCMPType                = "Digital's DDCMP protocol"
	CNICRawHDLCType              = "Raw HDLC"
	CNICRawIPType                = "Raw IP"
	CNICIPIPTunnelType           = "IPIP tunnel"
	CNICIP6IP6TunnelType         = "IP6IP6 tunnel"
	CNICFRADType                 = "Frame Relay Access Device"
	CNICSKIPType                 = "SKIP vif"
	CNICLoopbackType             = "Loopback device"
	CNICLocalTalkType            = "Localtalk device"
	CNICFDDIType                 = "Fiber Distributed Data Interface"
	CNICBIFType                  = "AP1000 BIF"
	CNICSITType                  = "sit0 device - IPv6-in-IPv4"
	CNICIPDDPType                = "IP over DDP tunneller"
	CNICIPGREType                = "GRE over IP"
	CNICPIMRegisterType          = "PIMSM register interface"
	CNICHPPIType                 = "High Performance Parallel Interface"
	CNICASHType                  = "Nexus 64Mbps Ash"
	CNICEcoNetType               = "Acorn Econet"
	CNICIRDAType                 = "Linux-IrDA"
	CNICFCPPType                 = "Point to point fibrechannel"
	CNICFCALType                 = "Fibrechannel arbitrated loop"
	CNICFCPLType                 = "Fibrechannel public loop"
	CNICFCFabricType             = "Fibrechannel fabric"
	CNICIEEE802TRType            = "Magic type ident for TR"
	CNICIEEE80211Type            = "IEEE 802.11"
	CNICIEEE802PrismType         = "IEEE 802.11 + Prism2 header"
	CNICIEEE802RadioTapType      = "IEEE 802.11 + radiotap header"
	CNICIEEE802154Type           = "IEEE 802.15.4"
	CNICIEEE802154MonitorType    = "IEEE 802.15.4 network monitor"
	CNICPhoNetType               = "PhoNet media type"
	CNICPhoNetPipeType           = "PhoNet pipe header"
	CNICCAIFType                 = "CAIF media type"
	CNICIP6GREType               = "GRE over IPv6"
	CNICNetlinkType              = "Netlink header"
	CNICIP6LoWPANType            = "IPv6 over LoWPAN"
	CNICVSockMonitorType         = "Vsock monitor header"
	CNICVoidType                 = "Void type, nothing is known"
	CNICNoneType                 = "zero header length"
)

var CTypes = map[uint16]Type{
	0:      CNICNETROMType,
	1:      CNICEthernetType,
	2:      CNICExperimentalEthernetType,
	3:      CNICAX25Type,
	4:      CNICPRONetType,
	5:      CNICChaosNetType,
	6:      CNICIEEE802Type,
	7:      CNICARCNetType,
	8:      CNICAppleTalkType,
	15:     CNICDLCIType,
	19:     CNICATMType,
	23:     CNICMetricomType,
	24:     CNICIEEE1394Type,
	27:     CNICEUI64Type,
	32:     CNICInfiniBandType,
	256:    CNICSlipType,
	257:    CNICCSlipType,
	258:    CNICSlip6Type,
	259:    CNICCSlip6Type,
	260:    CNICRSRVDType,
	264:    CNICAdaptType,
	270:    CNICRoseType,
	271:    CNICX25Type,
	272:    CNICHWX25Type,
	280:    CNICCANType,
	512:    CNICPPPType,
	513:    CNICCiscoHDLCType,
	516:    CNICLAPBType,
	517:    CNICDDCMPType,
	518:    CNICRawHDLCType,
	519:    CNICRawIPType,
	768:    CNICIPIPTunnelType,
	769:    CNICIP6IP6TunnelType,
	770:    CNICFRADType,
	771:    CNICSKIPType,
	772:    CNICLoopbackType,
	773:    CNICLocalTalkType,
	774:    CNICFDDIType,
	775:    CNICBIFType,
	776:    CNICSITType,
	777:    CNICIPDDPType,
	778:    CNICIPGREType,
	779:    CNICPIMRegisterType,
	780:    CNICHPPIType,
	781:    CNICASHType,
	782:    CNICEcoNetType,
	783:    CNICIRDAType,
	784:    CNICFCPPType,
	785:    CNICFCALType,
	786:    CNICFCPLType,
	787:    CNICFCFabricType,
	800:    CNICIEEE802TRType,
	801:    CNICIEEE80211Type,
	802:    CNICIEEE802PrismType,
	803:    CNICIEEE802RadioTapType,
	804:    CNICIEEE802154Type,
	805:    CNICIEEE802154MonitorType,
	820:    CNICPhoNetType,
	821:    CNICPhoNetPipeType,
	822:    CNICCAIFType,
	823:    CNICIP6GREType,
	824:    CNICNetlinkType,
	825:    CNICIP6LoWPANType,
	826:    CNICVSockMonitorType,
	0xffff: CNICVoidType,
	0xfffe: CNICNoneType,
}
