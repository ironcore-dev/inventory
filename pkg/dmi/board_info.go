package dmi

const (
	CUnsetBoardType                  = "Unset"
	CUnknownBoardType                = "Unknown"
	COtherBoardType                  = "Other"
	CServerBladeBoardType            = "Server Blade"
	CConnectivitySwitchBoardType     = "Connectivity Switch"
	CSystemManagementModuleBoardType = "System Management Module"
	CProcessorModuleBoardType        = "Processor Module"
	CIOModuleBoardType               = "I/O Module"
	CMemoryModuleBoardType           = "Memory Module"
	CDaughterBoardBoardType          = "Daughter board"
	CMotherboardBoardType            = "Motherboard (includes processor, memory, and I/O)"
	CProcessorMemoryModuleBoardType  = "Processor/Memory Module"
	CProcessorIOModuleBoardType      = "Processor/IO Module"
	CInterconnectBoardBoardType      = "Interconnect board"
)

var CBoardFeatureFlags = []string{
	"hosting board",
	"board requires at least one daughter board or auxiliary card to function properly",
	"board is removable; it is designed to be taken in and out of the chassis without impairing the function of the chassis",
	"board is replaceable; it is possible to replace (either as a field repair or as an upgrade) the board with a physically different board",
	"board is hot swappable; it is possible to replace the board with a physically different but equivalent board while power is applied to the board",
}

type BoardType string

var CBoardTypes = []BoardType{
	CUnsetBoardType,
	CUnknownBoardType,
	COtherBoardType,
	CServerBladeBoardType,
	CConnectivitySwitchBoardType,
	CSystemManagementModuleBoardType,
	CProcessorModuleBoardType,
	CIOModuleBoardType,
	CMemoryModuleBoardType,
	CDaughterBoardBoardType,
	CMotherboardBoardType,
	CProcessorMemoryModuleBoardType,
	CProcessorIOModuleBoardType,
	CInterconnectBoardBoardType,
}

type BoardInformationRefSpec struct {
	Manufacturer                   byte     `struc:"byte"`
	Product                        byte     `struc:"byte"`
	Version                        byte     `struc:"byte"`
	SerialNumber                   byte     `struc:"byte"`
	AssetTag                       byte     `struc:"byte"`
	FeatureFlags                   byte     `struc:"byte"`
	LocationInChassis              byte     `struc:"byte"`
	ChassisHandle                  uint16   `struc:"uint16"`
	Type                           byte     `struc:"byte"`
	NumberOfContainedObjectHandles byte     `struc:"byte,sizeof=ContainedObjectHandles"`
	ContainedObjectHandles         []uint16 `struc:"[]uint16"`
}

type BoardInformation struct {
	Manufacturer                   string
	Product                        string
	Version                        string
	SerialNumber                   string
	AssetTag                       string
	FeatureFlags                   []string
	LocationInChassis              string
	ChassisHandle                  uint16
	Type                           BoardType
	NumberOfContainedObjectHandles byte
	ContainedObjectHandles         []uint16
}

func BoardInformationFromSpec(ref *BoardInformationRefSpec, strings []string) *BoardInformation {
	info := &BoardInformation{}

	info.Manufacturer = emptyStringOrValue(ref.Manufacturer, strings)
	info.Product = emptyStringOrValue(ref.Product, strings)
	info.Version = emptyStringOrValue(ref.Version, strings)
	info.SerialNumber = emptyStringOrValue(ref.SerialNumber, strings)
	info.AssetTag = emptyStringOrValue(ref.AssetTag, strings)

	info.FeatureFlags = []string{}
	for i, ff := range CBoardFeatureFlags {
		idx := byte(1 << i)
		enabled := ref.FeatureFlags & idx
		if enabled != 0 {
			info.FeatureFlags = append(info.FeatureFlags, ff)
		}
	}

	info.LocationInChassis = emptyStringOrValue(ref.LocationInChassis, strings)
	info.ChassisHandle = ref.ChassisHandle

	if byte(len(CBoardTypes)) <= ref.Type {
		info.Type = CUnsetBoardType
	} else {
		info.Type = CBoardTypes[ref.Type]
	}

	info.NumberOfContainedObjectHandles = ref.NumberOfContainedObjectHandles
	info.ContainedObjectHandles = ref.ContainedObjectHandles

	return info
}
