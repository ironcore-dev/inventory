// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package smbiosinfo

import (
	"io"
	"os"
	"path/filepath"

	//nolint:revive
	_ "unsafe"

	"github.com/digitalocean/go-smbios/smbios"
	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/pkg/errors"
)

const (
	defaultRootPath             = "/"
	defaultMemPath              = "/dev/mem"
	defaultSysDMIPath           = "/sys/firmware/dmi/tables/DMI"
	defaultSysDMIEntryPointPath = "/sys/firmware/dmi/tables/smbios_entry_point"

	memSeekStartAddr = 0x000f0000
	memSeekEndAddr   = 0x000fffff
)

type SystemManager interface {
	GetUUID() (string, error)
	Close() error
}

//go:linkname sysfsStream github.com/digitalocean/go-smbios/smbios.sysfsStream
func sysfsStream(entryPoint, dmi string) (io.ReadCloser, smbios.EntryPoint, error)

//go:linkname memoryStream github.com/digitalocean/go-smbios/smbios.memoryStream
func memoryStream(rs io.ReadSeeker, startAddr, endAddr int) (io.ReadCloser, smbios.EntryPoint, error)

type systemManagement struct {
	readCloser io.ReadCloser
	entryPoint smbios.EntryPoint
	log        logger.Logger
}

func New(l logger.Logger) (SystemManager, error) {
	var stream io.ReadCloser
	var ep smbios.EntryPoint

	root := defaultRootPath
	if os.Getenv("ROOT") != "" {
		root = os.Getenv("ROOT")
	}

	sysDMIEntryPointPath := filepath.Join(root, defaultSysDMIEntryPointPath)
	_, err := os.Stat(filepath.Clean(sysDMIEntryPointPath))
	switch {
	case err == nil:
		sysDMIPath := filepath.Join(root, defaultSysDMIPath)
		stream, ep, err = sysfsStream(sysDMIEntryPointPath, filepath.Clean(sysDMIPath))
		if err != nil {
			return nil, errors.Wrap(err, "unable to access sysfs DMI stream")
		}
	case os.IsNotExist(err):
		devMemPath := filepath.Join(root, defaultMemPath)
		mem, openErr := os.Open(filepath.Clean(devMemPath))
		if openErr != nil {
			return nil, errors.Wrap(err, "unable to open /dev/mem")
		}
		defer func(mem *os.File) {
			if closeErr := mem.Close(); closeErr != nil {
				l.Info("can't close smbios memory folder", "error", closeErr)
			}
		}(mem)

		stream, ep, err = memoryStream(mem, memSeekStartAddr, memSeekEndAddr)
		if err != nil {
			return nil, errors.Wrap(err, "unable to access mem DMI stream")
		}
	default:
		return nil, errors.Wrapf(err, "unknown error while accessing DMI entry point")
	}

	return &systemManagement{
		readCloser: stream,
		entryPoint: ep,
		log:        l,
	}, nil
}

func (s *systemManagement) newDecoder() *smbios.Decoder {
	return smbios.NewDecoder(s.readCloser)
}

func (s *systemManagement) Close() error {
	return s.readCloser.Close()
}
