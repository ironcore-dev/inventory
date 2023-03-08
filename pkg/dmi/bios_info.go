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

package dmi

import "fmt"

var CCharacteristics = []string{
	"Reserved",
	"Reserved",
	"Unknown",
	"BIOS Characteristics are not supported",
	"ISA is supported",
	"MCA is supported",
	"EISA is supported",
	"PCI is supported",
	"PC card (PCMCIA) is supported",
	"Plug and Play is supported",
	"APM is supported",
	"BIOS is upgradeable (Flash)",
	"BIOS shadowing is allowed",
	"VL-VESA is supported",
	"ESCD support is available",
	"Boot from CD is supported",
	"Selectable boot is supported",
	"BIOS ROM is socketed(e.g. PLCC or SOP socket)",
	"Boot from PC card (PCMCIA) is supported",
	"EDD specification is supported",
	"Int 13h —Japanese floppy for NEC 9800 1.2MB(3.5”, 1Kbytes/sector, 360 RPM) is supported",
	"Int 13h—Japanese floppy for Toshiba 1.2MB(3.5”, 360 RPM) is supported",
	"Int 13h—5.25” / 360 KB floppy services are supported",
	"Int 13h—5.25” / 1.2MB floppy services are supported",
	"Int 13h—3.5” / 720 KB floppy services are supported",
	"Int 13h—3.5” / 2.88 MB floppy services are supported",
	"Int 5h, print screen Service is supported",
	"Int 9h, 8042 keyboard services are supported",
	"Int 14h, serial services are supported",
	"Int 17h, printer services are supported",
	"Int 10h, CGA/Mono Video Services are supported",
	"NEC PC-98",
}

var CCharacteristicsExtensionByte1 = []string{
	"ACPI is supported",
	"USB Legacy is supported",
	"AGP is supported",
	"I2O boot is supported",
	"LS-120 SuperDisk boot is supported",
	"ATAPI ZIP drive boot is supported",
	"1394 boot is supported",
	"Smart battery is supported",
}

var CCharacteristicsExtensionByte2 = []string{
	"BIOS Boot Specification is supported",
	"Function key-initiated network service boot is supported",
	"Enable targeted content distribution",
	"UEFI Specification is supported",
	"SMBIOS table describes a virtual machine",
}

var CCharacteristicsExtensions = [][]string{
	CCharacteristicsExtensionByte1,
	CCharacteristicsExtensionByte2,
}

type BIOSInformationRefSpec20 struct {
	Vendor                        byte   `struc:"byte"`
	Version                       byte   `struc:"byte"`
	StartingAddressSegment        uint16 `struc:"uint16"`
	ReleaseDate                   byte   `struc:"byte"`
	ROMSize                       byte   `struc:"byte"`
	Characteristics               uint64 `struc:"uint64,little"`
	CharacteristicsExtensions     []byte `struc:"[]byte,sizefrom=CharacteristicsExtensionsSize"`
	CharacteristicsExtensionsSize byte   `struc:"skip"`
}

type BIOSInformationRefSpec24 struct {
	BIOSInformationRefSpec20
	SystemMajorRelease                     byte `struc:"byte"`
	SystemMinorRelease                     byte `struc:"byte"`
	EmbeddedControllerFirmwareMajorRelease byte `struc:"byte"`
	EmbeddedControllerFirmwareMinorRelease byte `struc:"byte"`
}

type BIOSInformationRefSpec31 struct {
	BIOSInformationRefSpec24
	ExtendedROMSize uint16 `struc:"uint16"`
}

type BIOSInformation struct {
	Vendor                            string
	Version                           string
	StartingAddressSegment            string
	ReleaseDate                       string
	ROMSize                           uint64
	Characteristics                   []string
	SystemRelease                     string
	EmbeddedControllerFirmwareRelease string
}

func BIOSInformationFromSpec20(ref *BIOSInformationRefSpec20, strings []string) *BIOSInformation {
	info := &BIOSInformation{}

	// Reducing all values by one since structure contains element number
	// and we need an element index for the array
	info.Vendor = strings[ref.Vendor-1]
	info.Version = strings[ref.Version-1]
	info.StartingAddressSegment = fmt.Sprintf("%x", ref.StartingAddressSegment)
	info.ReleaseDate = strings[ref.ReleaseDate-1]
	info.ROMSize = (uint64(ref.ROMSize) + 1) * 64 * 1024
	info.Characteristics = make([]string, 0)

	for i, characteristic := range CCharacteristics {
		idx := uint64(1 << i)
		enabled := ref.Characteristics & idx
		if enabled != 0 {
			info.Characteristics = append(info.Characteristics, characteristic)
		}
	}

	for i, b := range ref.CharacteristicsExtensions {
		if i > len(CCharacteristicsExtensions) {
			break
		}
		for j, characteristic := range CCharacteristicsExtensions[i] {
			idx := uint8(1 << j)
			enabled := b & idx
			if enabled != 0 {
				info.Characteristics = append(info.Characteristics, characteristic)
			}
		}
	}

	return info
}

func BIOSInformationFromSpec24(ref *BIOSInformationRefSpec24, strings []string) *BIOSInformation {
	info := BIOSInformationFromSpec20(&ref.BIOSInformationRefSpec20, strings)

	info.SystemRelease = fmt.Sprintf("%d.%d", ref.SystemMajorRelease, ref.SystemMinorRelease)
	info.EmbeddedControllerFirmwareRelease = fmt.Sprintf("%d.%d", ref.EmbeddedControllerFirmwareMajorRelease, ref.EmbeddedControllerFirmwareMinorRelease)

	return info
}

func BIOSInformationFromSpec31(ref *BIOSInformationRefSpec31, strings []string) *BIOSInformation {
	info := BIOSInformationFromSpec24(&ref.BIOSInformationRefSpec24, strings)

	if ref.ExtendedROMSize == 0 {
		return info
	}

	// (11)(11111111111111)
	// first group (2 bits) - measurement unit
	// second group (14 bits) - size
	unit := ref.ExtendedROMSize >> 14
	size := ref.ExtendedROMSize << 2 >> 2

	switch unit {
	// MB
	case 0x00:
		info.ROMSize = uint64(size) * 1024 * 1024
	// GB
	case 0x01:
		info.ROMSize = uint64(size) * 1024 * 1024 * 1024
	}

	return info

}
