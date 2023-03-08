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

import (
	"bytes"

	"github.com/digitalocean/go-smbios/smbios"
	"github.com/lunixbochs/struc"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

type Svc struct {
	printer   *printer.Svc
	RawDMISvc *RawSvc
}

func NewSvc(printer *printer.Svc, rawDMISvc *RawSvc) *Svc {
	return &Svc{
		printer:   printer,
		RawDMISvc: rawDMISvc,
	}
}

func (s *Svc) GetData() (*DMI, error) {
	rawDmi, err := s.RawDMISvc.GetRaw()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get SMBIOS stream")
	}

	defer rawDmi.Stream.Close()

	decoder := smbios.NewDecoder(rawDmi.Stream)
	structures, err := decoder.Decode()
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode SMBIOS stream")
	}

	version := NewSMBIOSVersion(rawDmi.EntryPoint.Version())

	dmi := &DMI{
		Version: version,
	}

	for _, structure := range structures {
		switch structure.Header.Type {
		case CBoardInformationHeaderType:
			boardInfo, err := s.parseBoardInformation(structure, version)
			if err != nil {
				s.printer.VErr(errors.Wrap(err, "unable to parse board info"))
				continue
			}
			if dmi.BoardInformation == nil {
				dmi.BoardInformation = []BoardInformation{}
			}
			dmi.BoardInformation = append(dmi.BoardInformation, *boardInfo)
		case CBIOSInformationHeaderType:
			biosInfo, err := s.parseBiosInformation(structure, version)
			if err != nil {
				s.printer.VErr(errors.Wrap(err, "unable to parse BIOS info"))
			}
			dmi.BIOSInformation = biosInfo
		case CSystemInformationHeaderType:
			systemInfo, err := s.parseSystemInformation(structure, version)
			if err != nil {
				s.printer.VErr(errors.Wrap(err, "unable to parse system info"))
			}
			dmi.SystemInformation = systemInfo
		}
	}

	return dmi, nil
}

func (s *Svc) parseBoardInformation(structure *smbios.Structure, version *SMBIOSVersion) (*BoardInformation, error) {
	ref := &BoardInformationRefSpec{}

	if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
		return nil, errors.Wrap(err, "unable to unpack structure")
	}

	return BoardInformationFromSpec(ref, structure.Strings), nil
}

func (s *Svc) parseBiosInformation(structure *smbios.Structure, version *SMBIOSVersion) (*BIOSInformation, error) {
	// Spec contains info only for 2.0+
	if version.Lesser(&SMBIOSVersion{2, 0, 0}) {
		return &BIOSInformation{}, nil
	}

	// 3.1+
	if version.GreaterOrEqual(&SMBIOSVersion{3, 1, 0}) {
		ref := &BIOSInformationRefSpec31{}
		// Subtracting base structure offset (12b)
		// Subtracting 4 bytes going after extension byte array
		// Subtracting 2 bytes going after extension byte array
		ref.CharacteristicsExtensionsSize = structure.Header.Length - 0x12 - 0x4 - 0x2
		if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
			return nil, errors.Wrap(err, "unable to unpack structure")
		}
		return BIOSInformationFromSpec31(ref, structure.Strings), nil
	}

	// 2.4+
	if version.GreaterOrEqual(&SMBIOSVersion{2, 4, 0}) {
		ref := &BIOSInformationRefSpec24{}
		// Subtracting base structure offset (12b)
		// Subtracting 4 bytes going after extension byte array
		ref.CharacteristicsExtensionsSize = structure.Header.Length - 0x12 - 0x4
		if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
			return nil, errors.Wrap(err, "unable to unpack structure")
		}
		return BIOSInformationFromSpec24(ref, structure.Strings), nil
	}

	// 2.0+
	ref := &BIOSInformationRefSpec20{}
	// Subtracting base structure offset (12b)
	ref.CharacteristicsExtensionsSize = structure.Header.Length - 0x12
	if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
		return nil, errors.Wrap(err, "unable to unpack structure")
	}
	return BIOSInformationFromSpec20(ref, structure.Strings), nil
}

func (s *Svc) parseSystemInformation(structure *smbios.Structure, version *SMBIOSVersion) (*SystemInformation, error) {
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
