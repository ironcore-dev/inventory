package nic

import (
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	CNICDevicePCIAddressPath = "/device"

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
)

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
	AddressLength       int64
	Broadcast           string
	Carrier             bool
	CarrierChanges      int64
	CarrierDownCount    int64
	CarrierUpCount      int64
	DevID               string
	DevPort             uint8
	Dormant             bool
	Duplex              string
	Flags               *Flags
	InterfaceAlias      string
	InterfaceIndex      int64
	InterfaceLink       int64
	LinkMode            LinkMode
	MTU                 uint16
	NameAssignType      NameAssignType
	NetDevGroup         int64
	OperationalState    string
	PhysicalPortID      string
	PhysicalPortName    string
	PhysicalSwitchID    string
	Speed               uint32
	Testing             bool
	TransmitQueueLength int64
	Type                Type
	Lanes               uint8
	FEC                 string
}

func (n *Device) defPCIAddress(thePath string) error {
	filePath := path.Join(thePath, CNICDevicePCIAddressPath)

	if _, err := os.Stat(filePath); err == nil {
		linkPath, err := filepath.EvalSymlinks(filePath)
		if err != nil {
			return errors.Wrapf(err, "unable resolve symlink with path %s", filePath)
		}
		fileInfo, err := os.Stat(linkPath)
		if err != nil {
			return errors.Wrapf(err, "unable to get stat for path %s", filePath)
		}

		n.PCIAddress = fileInfo.Name()
	}

	return nil
}
