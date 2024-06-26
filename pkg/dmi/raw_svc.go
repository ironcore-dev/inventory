// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package dmi

import (
	"io"
	"os"
	"path"
	_ "unsafe"

	"github.com/digitalocean/go-smbios/smbios"
	"github.com/pkg/errors"
)

const (
	CDevMemPath           = "/dev/mem"
	CSysDMIPath           = "/sys/firmware/dmi/tables/DMI"
	CSysDMIEntryPointPath = "/sys/firmware/dmi/tables/smbios_entry_point"

	CMemSeekStartAddr = 0x000f0000
	CMemSeekEndAddr   = 0x000fffff
)

//go:linkname sysfsStream github.com/digitalocean/go-smbios/smbios.sysfsStream
func sysfsStream(entryPoint, dmi string) (io.ReadCloser, smbios.EntryPoint, error)

//go:linkname memoryStream github.com/digitalocean/go-smbios/smbios.memoryStream
func memoryStream(rs io.ReadSeeker, startAddr, endAddr int) (io.ReadCloser, smbios.EntryPoint, error)

type RawSvc struct {
	devMemPath           string
	sysDMIPath           string
	sysDMIEntryPointPath string
}

func NewRawSvc(basePath string) *RawSvc {
	return &RawSvc{
		devMemPath:           path.Join(basePath, CDevMemPath),
		sysDMIPath:           path.Join(basePath, CSysDMIPath),
		sysDMIEntryPointPath: path.Join(basePath, CSysDMIEntryPointPath),
	}
}

func (s *RawSvc) GetRaw() (*Raw, error) {
	var stream io.ReadCloser
	var ep smbios.EntryPoint

	_, err := os.Stat(s.sysDMIEntryPointPath)
	switch {
	case err == nil:
		stream, ep, err = sysfsStream(s.sysDMIEntryPointPath, s.sysDMIPath)
		if err != nil {
			return nil, errors.Wrap(err, "unable to access sysfs DMI stream")
		}
	case os.IsNotExist(err):
		mem, err := os.Open(s.devMemPath)
		if err != nil {
			return nil, errors.Wrap(err, "unable to open /dev/mem")
		}
		defer mem.Close()

		stream, ep, err = memoryStream(mem, CMemSeekStartAddr, CMemSeekEndAddr)
		if err != nil {
			return nil, errors.Wrap(err, "unable to access mem DMI stream")
		}
	default:
		return nil, errors.Wrapf(err, "unknown error while accessing DMI entry point")
	}

	return &Raw{
		Stream:     stream,
		EntryPoint: ep,
	}, nil
}
