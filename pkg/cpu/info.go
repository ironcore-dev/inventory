// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package cpu

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	CCPUInfoProcessorKey       = "processor"
	CCPUInfoVendorIDKey        = "vendor_id"
	CCPUInfoCPUFamilyKey       = "cpu family"
	CCPUInfoModelKey           = "model"
	CCPUInfoModelNameKey       = "model name"
	CCPUInfoSteppingKey        = "stepping"
	CCPUInfoMicrocodeKey       = "microcode"
	CCPUInfoCPUMHzKey          = "cpu MHz"
	CCPUInfoCacheSizeKey       = "cache size"
	CCPUInfoPhysicalIDKey      = "physical id"
	CCPUInfoSiblingsKey        = "siblings"
	CCPUInfoCoreIdKey          = "core id"
	CCPUInfoCpuCoresKey        = "cpu cores"
	CCPUInfoAPICIDKey          = "apicid"
	CCPUInfoInitialAPICIDKey   = "initial apicid"
	CCPUInfoFPUKey             = "fpu"
	CCPUInfoFPUExceptionKey    = "fpu_exception"
	CCPUInfoCPUIDLevelKey      = "cpuid level"
	CCPUInfoWPKey              = "wp"
	CCPUInfoFlagsKey           = "flags"
	CCPUInfoVMXFlagsKey        = "vmx flags"
	CCPUInfoBugsKey            = "bugs"
	CCPUInfoBogoMIPSKey        = "bogomips"
	CCPUInfoCLFlushSizeKey     = "clflush size"
	CCPUInfoCacheAlignmentKey  = "cache_alignment"
	CCPUInfoAddressSizesKey    = "address sizes"
	CCPUInfoPowerManagementKey = "power management"
)

type Info struct {
	Processor       uint64
	VendorID        string
	CPUFamily       string
	Model           string
	ModelName       string
	Stepping        string
	Microcode       string
	CPUMHz          string
	CacheSize       string
	PhysicalID      uint64
	Siblings        uint64
	CoreID          string
	CpuCores        uint64
	APICID          string
	InitialAPICID   string
	FPU             bool
	FPUException    bool
	CPUIDLevel      uint64
	WP              bool
	Flags           []string
	VMXFlags        []string
	Bugs            []string
	BogoMIPS        string
	CLFlushSize     uint64
	CacheAlignment  uint64
	AddressSizes    string
	PowerManagement string
}

func (ci *Info) setField(key string, val string) error {
	switch key {
	case CCPUInfoProcessorKey:
		v, err := strconv.ParseUint(val, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "unable to convert %s to uint", val)
		}
		ci.Processor = v
	case CCPUInfoVendorIDKey:
		ci.VendorID = val
	case CCPUInfoCPUFamilyKey:
		ci.CPUFamily = val
	case CCPUInfoModelKey:
		ci.Model = val
	case CCPUInfoModelNameKey:
		ci.ModelName = val
	case CCPUInfoSteppingKey:
		ci.Stepping = val
	case CCPUInfoMicrocodeKey:
		ci.Microcode = val
	case CCPUInfoCPUMHzKey:
		ci.CPUMHz = val
	case CCPUInfoCacheSizeKey:
		ci.CacheSize = val
	case CCPUInfoPhysicalIDKey:
		v, err := strconv.ParseUint(val, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "unable to convert %s to uint", val)
		}
		ci.PhysicalID = v
	case CCPUInfoSiblingsKey:
		v, err := strconv.ParseUint(val, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "unable to convert %s to uint", val)
		}
		ci.Siblings = v
	case CCPUInfoCoreIdKey:
		ci.CoreID = val
	case CCPUInfoCpuCoresKey:
		v, err := strconv.ParseUint(val, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "unable to convert %s to uint", val)
		}
		ci.CpuCores = v
	case CCPUInfoAPICIDKey:
		ci.APICID = val
	case CCPUInfoInitialAPICIDKey:
		ci.InitialAPICID = val
	case CCPUInfoFPUKey:
		var v bool
		if val == "yes" {
			v = true
		}
		ci.FPU = v
	case CCPUInfoFPUExceptionKey:
		var v bool
		if val == "yes" {
			v = true
		}
		ci.FPUException = v
	case CCPUInfoCPUIDLevelKey:
		v, err := strconv.ParseUint(val, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "unable to convert %s to uint", val)
		}
		ci.CPUIDLevel = v
	case CCPUInfoWPKey:
		var v bool
		if val == "yes" {
			v = true
		}
		ci.WP = v
	case CCPUInfoFlagsKey:
		v := strings.Split(val, " ")
		ci.Flags = v
	case CCPUInfoVMXFlagsKey:
		v := strings.Split(val, " ")
		ci.VMXFlags = v
	case CCPUInfoBugsKey:
		v := strings.Split(val, " ")
		ci.Bugs = v
	case CCPUInfoBogoMIPSKey:
		ci.BogoMIPS = val
	case CCPUInfoCLFlushSizeKey:
		v, err := strconv.ParseUint(val, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "unable to convert %s to uint", val)
		}
		ci.CLFlushSize = v
	case CCPUInfoCacheAlignmentKey:
		v, err := strconv.ParseUint(val, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "unable to convert %s to uint", val)
		}
		ci.CacheAlignment = v
	case CCPUInfoAddressSizesKey:
		ci.AddressSizes = val
	case CCPUInfoPowerManagementKey:
		ci.PowerManagement = val
	default:
		return errors.Errorf("unknown key %s from cpuinfo", key)
	}
	return nil
}
