// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package block

import (
	"path"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/file"
	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CQueueRotationalPath        = "/queue/rotational"
	CQueuePhysicalBlockSizePath = "/queue/physical_block_size"
	CQueueLogicalBlockSizePath  = "/queue/logical_block_size"
	CQueueHWSectorSizePath      = "/queue/hw_sector_size"
	CDeviceVendorPath           = "/device/vendor"
	CDeviceModelPath            = "/device/model"
	CDeviceSerialPath           = "/device/serial"
	CDeviceNumaNodePath         = "/device/numa_node"
	CDeviceFirmwareRevPath      = "/device/firmware_rev"
	CDeviceState                = "/device/state"
	CWWIDPath                   = "/wwid"
	CRemovablePath              = "/removable"
	CSizePath                   = "/size"
	CReadOnlyPath               = "/ro"
)

type DeviceSvc struct {
	printer            *printer.Svc
	partitionTableSvc  *PartitionTableSvc
	blockDeviceStatSvc *DeviceStatSvc
}

func NewDeviceSvc(printer *printer.Svc, partTableSvc *PartitionTableSvc, blockDeviceStatSvc *DeviceStatSvc) *DeviceSvc {
	return &DeviceSvc{
		printer:            printer,
		partitionTableSvc:  partTableSvc,
		blockDeviceStatSvc: blockDeviceStatSvc,
	}
}

func (s *DeviceSvc) GetDevice(thePath string, name string) (*Device, error) {
	device := &Device{
		path: thePath,
		Name: name,
	}

	defs := []func(*Device) error{
		s.setType,
		s.setRotational,
		s.setRemovable,
		s.setReadOnly,
		s.setVendor,
		s.setModel,
		s.setSerial,
		s.setWWID,
		s.setFirmwareRev,
		s.setState,
		s.setHWSectorSize,
		s.setPhysicalBlockSize,
		s.setLogicalBlockSize,
		s.setSize,
		s.setNumaNodeID,
		s.setPartitionTable,
		s.setStat,
	}

	for _, def := range defs {
		err := def(device)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to set block device property"))
		}
	}

	return device, nil
}

func (s *DeviceSvc) setType(bd *Device) error {
	for k, v := range CDiskMap {
		if v.MatchString(bd.Name) {
			bd.Type = k
			break
		}
	}

	if bd.Type == "" {
		return errors.Errorf("unhandled block device type %s found", bd.Name)
	}

	return nil
}

func (s *DeviceSvc) setRotational(bd *Device) error {
	rotationalPath := path.Join(bd.path, CQueueRotationalPath)
	rotational, err := file.ToBool(rotationalPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", rotationalPath)
	}

	bd.Rotational = rotational

	return nil
}

func (s *DeviceSvc) setVendor(bd *Device) error {
	vendorPath := path.Join(bd.path, CDeviceVendorPath)
	vendor, err := file.ToString(vendorPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", vendorPath)
	}

	bd.Vendor = vendor

	return nil
}

func (s *DeviceSvc) setModel(bd *Device) error {
	modelPath := path.Join(bd.path, CDeviceModelPath)
	model, err := file.ToString(modelPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", modelPath)
	}

	bd.Model = model

	return nil
}

func (s *DeviceSvc) setWWID(bd *Device) error {
	wwidPath := path.Join(bd.path, CWWIDPath)
	wwid, err := file.ToString(wwidPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", wwidPath)
	}

	bd.WWID = wwid

	return nil
}

func (s *DeviceSvc) setRemovable(bd *Device) error {
	removablePath := path.Join(bd.path, CRemovablePath)
	removable, err := file.ToBool(removablePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", removablePath)
	}

	bd.Removable = removable

	return nil
}

func (s *DeviceSvc) setSerial(bd *Device) error {
	serialPath := path.Join(bd.path, CDeviceSerialPath)
	serial, err := file.ToString(serialPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", serialPath)
	}

	bd.Serial = serial

	return nil
}

func (s *DeviceSvc) setSize(bd *Device) error {
	sizePath := path.Join(bd.path, CSizePath)
	size, err := file.ToUint64(sizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", sizePath)
	}

	bd.Size = size * bd.HWSectorSize

	return nil
}

func (s *DeviceSvc) setPhysicalBlockSize(bd *Device) error {
	blockSizePath := path.Join(bd.path, CQueuePhysicalBlockSizePath)
	blockSize, err := file.ToUint64(blockSizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", blockSizePath)
	}

	bd.PhysicalBlockSize = blockSize

	return nil
}

func (s *DeviceSvc) setLogicalBlockSize(bd *Device) error {
	blockSizePath := path.Join(bd.path, CQueueLogicalBlockSizePath)
	blockSize, err := file.ToUint64(blockSizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", blockSizePath)
	}

	bd.LogicalBlockSize = blockSize

	return nil
}

func (s *DeviceSvc) setHWSectorSize(bd *Device) error {
	sectorSizePath := path.Join(bd.path, CQueueHWSectorSizePath)
	sectorSize, err := file.ToUint64(sectorSizePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", sectorSizePath)
	}

	bd.HWSectorSize = sectorSize

	return nil
}

func (s *DeviceSvc) setNumaNodeID(bd *Device) error {
	numaPath := path.Join(bd.path, CDeviceNumaNodePath)
	numa, err := file.ToUint64(numaPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", numaPath)
	}

	bd.NUMANodeID = numa

	return nil
}

func (s *DeviceSvc) setReadOnly(bd *Device) error {
	roPath := path.Join(bd.path, CReadOnlyPath)
	ro, err := file.ToBool(roPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", roPath)
	}

	bd.ReadOnly = ro

	return nil
}

func (s *DeviceSvc) setFirmwareRev(bd *Device) error {
	fwPath := path.Join(bd.path, CDeviceFirmwareRevPath)
	fw, err := file.ToString(fwPath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", fwPath)
	}

	bd.FirmwareRevision = fw

	return nil
}

func (s *DeviceSvc) setState(bd *Device) error {
	statePath := path.Join(bd.path, CDeviceState)
	state, err := file.ToString(statePath)
	if err != nil {
		return errors.Wrapf(err, "unable to get value from file %s", statePath)
	}

	bd.State = state

	return nil
}

func (s *DeviceSvc) setPartitionTable(bd *Device) error {
	table, err := s.partitionTableSvc.GetPartitionTable(bd.Name)
	if err != nil {
		return errors.Wrap(err, "unable to get partition table")
	}

	bd.PartitionTable = table

	return nil
}

func (s *DeviceSvc) setStat(bd *Device) error {
	stat, err := s.blockDeviceStatSvc.GetDeviceStat(bd.path)
	if err != nil {
		return errors.Wrap(err, "unable to collect stats")
	}

	bd.Stat = stat

	return nil
}
