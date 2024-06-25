// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package pci

type Bus struct {
	ID      string
	Devices []Device
}

func NewBus(id string, devices []Device) *Bus {
	return &Bus{
		ID:      id,
		Devices: devices,
	}
}
