// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package smbiosinfo

import (
	"bytes"
	"fmt"

	"github.com/digitalocean/go-smbios/smbios"
	"github.com/lunixbochs/struc"
	"github.com/pkg/errors"

	bencherr "github.com/onmetal/inventory/internal/errors"
)

type WakeUpType string

const (
	CSystemInformationHeaderType = 1
)

type SystemInformationRefSpec20 struct {
	Manufacturer byte `struc:"byte"`
	ProductName  byte `struc:"byte"`
	Version      byte `struc:"byte"`
	SerialNumber byte `struc:"byte"`
}

type SystemInformationRefSpec21 struct { //nolint:govet //reason: unexpected behavior when sorted
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

func (s *systemManagement) GetUUID() (string, error) {
	version := NewSMBIOSVersion(s.entryPoint.Version())

	d := s.newDecoder()
	structures, err := d.Decode()
	if err != nil {
		return "", bencherr.Unknown(err.Error())
	}

	for _, structure := range structures {
		if structure.Header.Type != CSystemInformationHeaderType {
			continue
		}
		systemInfo, err := parseSystemInformation(structure, version)
		if err != nil {
			return "", err
		}
		return systemInfo.UUID, nil
	}
	return "", bencherr.NotFound("uuid")
}

func parseSystemInformation(structure *smbios.Structure, version *SMBIOSVersion) (*SystemInformation, error) {
	// Spec contains info only for 2.0+
	if version.Lesser(&SMBIOSVersion{2, 0, 0}) {
		return &SystemInformation{}, nil
	}

	// 2.4+
	if version.GreaterOrEqual(&SMBIOSVersion{2, 4, 0}) {
		ref := &SystemInformationRefSpec24{}
		if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
			return nil, errors.Wrap(err, "unable to unpack structure")
		}
		return SystemInformationFromSpec24(ref, structure.Strings), nil
	}

	// 2.1+
	if version.GreaterOrEqual(&SMBIOSVersion{2, 1, 0}) {
		ref := &SystemInformationRefSpec21{}
		if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
			return nil, errors.Wrap(err, "unable to unpack structure")
		}
		return SystemInformationFromSpec21(ref, structure.Strings), nil
	}

	// 2.0+
	ref := &SystemInformationRefSpec20{}
	if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
		return nil, errors.Wrap(err, "unable to unpack structure")
	}
	return SystemInformationFromSpec20(ref, structure.Strings), nil
}

func SystemInformationFromSpec24(ref *SystemInformationRefSpec24, strings []string) *SystemInformation {
	info := SystemInformationFromSpec21(&ref.SystemInformationRefSpec21, strings)

	// Reducing all values by one since structure contains element number
	// and we need an element index for the array
	info.SKUNumber = strings[ref.SKUNumber-1]
	info.Family = strings[ref.Family-1]

	return info
}

func SystemInformationFromSpec21(ref *SystemInformationRefSpec21, strings []string) *SystemInformation {
	info := SystemInformationFromSpec20(&ref.SystemInformationRefSpec20, strings)
	info.UUID = calculateUUID(ref)
	return info
}

func calculateUUID(ref *SystemInformationRefSpec21) string {
	// According to SMBIOS spec UUID bytes should be ordered as
	// 33 22 11 00 55 44 77 66 88 99 AA BB CC DD EE FF
	// see 7.2.1 System — UUID
	uuidBytes := make([]byte, len(ref.UUID))
	copy(uuidBytes, ref.UUID)
	swapBytesInSlice(uuidBytes, 0, 3) //nolint:gomnd //Reason: false-positive
	swapBytesInSlice(uuidBytes, 1, 2) //nolint:gomnd //Reason: false-positive
	swapBytesInSlice(uuidBytes, 4, 5) //nolint:gomnd //Reason: false-positive
	swapBytesInSlice(uuidBytes, 6, 7) //nolint:gomnd //Reason: false-positive

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:])
}

func swapBytesInSlice(slice []byte, a, b int) {
	slice[a], slice[b] = slice[b], slice[a]
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

func emptyStringOrValue(index byte, strings []string) string {
	if index == byte(0) || int(index) > len(strings) {
		str := ""
		return str
	}
	return strings[index-1]
}
