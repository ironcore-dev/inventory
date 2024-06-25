// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package dmi

import (
	"fmt"
)

type WakeUpType string

const (
	CReservedWakeUpType        WakeUpType = "Reserved"
	COtherWakeUpType           WakeUpType = "Other"
	CUnknownWakeUpType         WakeUpType = "Unknown"
	CAPMTimerWakeUpType        WakeUpType = "APM Timer"
	CModemRingWakeUpType       WakeUpType = "Modem Ring"
	CLANRemoteWakeUpType       WakeUpType = "LAN Remote"
	CPowerSwitchWakeUpType     WakeUpType = "Power Switch"
	CPCIPMEWakeUpType          WakeUpType = "PCI PME#"
	CACPowerRestoredWakeUpType WakeUpType = "AC Power Restored"
)

var wakeUpTypes = map[byte]WakeUpType{
	0: CReservedWakeUpType,
	1: COtherWakeUpType,
	2: CUnknownWakeUpType,
	3: CAPMTimerWakeUpType,
	4: CModemRingWakeUpType,
	5: CLANRemoteWakeUpType,
	6: CPowerSwitchWakeUpType,
	7: CPCIPMEWakeUpType,
	8: CACPowerRestoredWakeUpType,
}

type SystemInformationRefSpec20 struct {
	Manufacturer byte `struc:"byte"`
	ProductName  byte `struc:"byte"`
	Version      byte `struc:"byte"`
	SerialNumber byte `struc:"byte"`
}

type SystemInformationRefSpec21 struct {
	SystemInformationRefSpec20
	UUID       []byte `struc:"[16]byte"`
	WakeUpType byte   `struc:"byte"`
}

type SystemInformationRefSpec24 struct {
	SystemInformationRefSpec21
	SKUNumber byte `struc:"byte"`
	Family    byte `struc:"byte"`
}

type SystemInformation struct {
	Manufacturer string
	ProductName  string
	Version      string
	SerialNumber string
	UUID         string
	WakeUpType   WakeUpType
	SKUNumber    string
	Family       string
}

func SystemInformationFromSpec20(ref *SystemInformationRefSpec20, strings []string) *SystemInformation {
	// Reducing all values by one since structure contains element number
	// and we need an element index for the array
	info := &SystemInformation{
		Manufacturer: emptyStringOrValue(ref.Manufacturer, strings),
		ProductName:  emptyStringOrValue(ref.ProductName, strings),
		Version:      emptyStringOrValue(ref.Version, strings),
		SerialNumber: emptyStringOrValue(ref.SerialNumber, strings),
	}

	return info
}

func SystemInformationFromSpec21(ref *SystemInformationRefSpec21, strings []string) *SystemInformation {
	info := SystemInformationFromSpec20(&ref.SystemInformationRefSpec20, strings)

	// According to SMBIOS spec UUID bytes should be ordered as
	// 33 22 11 00 55 44 77 66 88 99 AA BB CC DD EE FF
	// see 7.2.1 System — UUID
	uuidBytes := make([]byte, len(ref.UUID))
	copy(uuidBytes, ref.UUID)
	swapBytesInSlice(uuidBytes, 0, 3)
	swapBytesInSlice(uuidBytes, 1, 2)
	swapBytesInSlice(uuidBytes, 4, 5)
	swapBytesInSlice(uuidBytes, 6, 7)

	info.UUID = fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:])
	info.WakeUpType = wakeUpTypes[ref.WakeUpType]

	return info
}

func SystemInformationFromSpec24(ref *SystemInformationRefSpec24, strings []string) *SystemInformation {
	info := SystemInformationFromSpec21(&ref.SystemInformationRefSpec21, strings)

	// Reducing all values by one since structure contains element number
	// and we need an element index for the array
	info.SKUNumber = strings[ref.SKUNumber-1]
	info.Family = strings[ref.Family-1]

	return info
}

func swapBytesInSlice(slice []byte, a int, b int) {
	tmp := slice[a]
	slice[a] = slice[b]
	slice[b] = tmp
}

func emptyStringOrValue(index byte, strings []string) string {
	if index == byte(0) || int(index) > len(strings) {
		str := ""
		return str
	} else {
		return strings[index-1]
	}
}
