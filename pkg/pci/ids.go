// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package pci

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	CResPCIIDsFilePath = "./res/pci.ids"

	CCommentLinePrefix = "#"

	CVendorLinePattern    = "^([[:xdigit:]]{4})\\s+(.*)$"
	CDeviceLinePattern    = "^\t([[:xdigit:]]{4})\\s+(.*)$"
	CSubsystemLinePattern = "^\t\t([[:xdigit:]]{4})\\s+([[:xdigit:]]{4})\\s+(.*)$"

	CClassLinePattern                = "^C ([[:xdigit:]]{2})\\s+(.*)$"
	CSubclassLinePattern             = "^\t([[:xdigit:]]{2})\\s+(.*)$"
	CProgrammingInterfaceLinePattern = "^\t\t([[:xdigit:]]{2})\\s+(.*)$"
)

var CVendorLineRegexp = regexp.MustCompile(CVendorLinePattern)
var CDeviceLineRegexp = regexp.MustCompile(CDeviceLinePattern)
var CSubsystemLineRegexp = regexp.MustCompile(CSubsystemLinePattern)

var CClassLineRegexp = regexp.MustCompile(CClassLinePattern)
var CSubclassLineRegexp = regexp.MustCompile(CSubclassLinePattern)
var CProgrammingInterfaceLineRegexp = regexp.MustCompile(CProgrammingInterfaceLinePattern)

type Vendor struct {
	ID      string
	Name    string
	Devices map[string]Type
}

type Type struct {
	ID         string
	Name       string
	Subsystems map[string]Subsystem
}

type Subsystem struct {
	SubvendorID string
	SubdeviceID string
	Name        string
}

type Class struct {
	ID         string
	Name       string
	Subclasses map[string]Subclass
}

type Subclass struct {
	ID                    string
	Name                  string
	ProgrammingInterfaces map[string]ProgrammingInterface
}

type ProgrammingInterface struct {
	ID   string
	Name string
}

type IDs struct {
	Vendors map[string]Vendor
	Classes map[string]Class
}

func NewIDs() (*IDs, error) {
	procPath, err := os.Executable()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get executable path")
	}
	procDir := filepath.Dir(procPath)

	idsPath := path.Join(procDir, CResPCIIDsFilePath)
	fh, err := os.Open(idsPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open file %s", CResPCIIDsFilePath)
	}

	defer fh.Close()

	fs := bufio.NewScanner(fh)

	vendors := map[string]Vendor{}
	classes := map[string]Class{}

	var vendor *Vendor
	var class *Class
	var device *Type
	var subclass *Subclass

	for fs.Scan() {
		line := fs.Text()

		// skip empty line
		if strings.TrimSpace(line) == "" {
			continue
		}
		// skip comment line
		if strings.HasPrefix(line, CCommentLinePrefix) {
			continue
		}

		groups := CVendorLineRegexp.FindStringSubmatch(line)
		if len(groups) == 3 {
			if device != nil {
				vendor.Devices[device.ID] = *device
				device = nil
			}
			if vendor != nil {
				vendors[vendor.ID] = *vendor
			}

			vendor = &Vendor{
				ID:      groups[1],
				Name:    groups[2],
				Devices: map[string]Type{},
			}
			continue
		}

		groups = CDeviceLineRegexp.FindStringSubmatch(line)
		if len(groups) == 3 {
			if device != nil {
				vendor.Devices[device.ID] = *device
			}

			device = &Type{
				ID:         groups[1],
				Name:       groups[2],
				Subsystems: map[string]Subsystem{},
			}
			continue
		}

		groups = CSubsystemLineRegexp.FindStringSubmatch(line)
		if len(groups) == 4 {
			subsystem := Subsystem{
				SubvendorID: groups[1],
				SubdeviceID: groups[2],
				Name:        groups[3],
			}

			device.Subsystems[subsystem.SubvendorID+subsystem.SubdeviceID] = subsystem
			continue
		}

		groups = CClassLineRegexp.FindStringSubmatch(line)
		if len(groups) == 3 {
			if subclass != nil {
				class.Subclasses[subclass.ID] = *subclass
				subclass = nil
			}
			if class != nil {
				classes[class.ID] = *class
			}

			class = &Class{
				ID:         groups[1],
				Name:       groups[2],
				Subclasses: map[string]Subclass{},
			}
			continue
		}

		groups = CSubclassLineRegexp.FindStringSubmatch(line)
		if len(groups) == 3 {
			if subclass != nil {
				class.Subclasses[subclass.ID] = *subclass
			}

			subclass = &Subclass{
				ID:                    groups[1],
				Name:                  groups[2],
				ProgrammingInterfaces: map[string]ProgrammingInterface{},
			}
			continue
		}

		groups = CProgrammingInterfaceLineRegexp.FindStringSubmatch(line)
		if len(groups) == 3 {
			progIface := ProgrammingInterface{
				ID:   groups[1],
				Name: groups[2],
			}
			subclass.ProgrammingInterfaces[progIface.ID] = progIface
			continue
		}
	}

	if device != nil {
		vendor.Devices[device.ID] = *device
	}

	if vendor != nil {
		vendors[vendor.ID] = *vendor
	}

	if subclass != nil {
		class.Subclasses[subclass.ID] = *subclass
	}

	if class != nil {
		classes[class.ID] = *class
	}

	return &IDs{
		Vendors: vendors,
		Classes: classes,
	}, err
}
