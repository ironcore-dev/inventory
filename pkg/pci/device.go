// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package pci

type DeviceType struct {
	ID   string
	Name string
}

type DeviceVendor struct {
	ID   string
	Name string
}

type DeviceSubtype struct {
	ID   string
	Name string
}

type DeviceClass struct {
	ID   string
	Name string
}

type DeviceSubclass struct {
	ID   string
	Name string
}

type DeviceProgrammingInterface struct {
	ID   string
	Name string
}

type Device struct {
	Address              string
	Vendor               *DeviceVendor
	Type                 *DeviceType
	Subvendor            *DeviceVendor
	Subtype              *DeviceSubtype
	Class                *DeviceClass
	Subclass             *DeviceSubclass
	ProgrammingInterface *DeviceProgrammingInterface
}
