// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package block

import "regexp"

const (
	CFloppyDiskPattern = "fd(\\d+)"
	CCDROMDiskPattern  = "(sr|scd)(\\d+)"
	CSCSIDiskPattern   = "sd(\\w+)"
	CIDEDiskPattern    = "hd(\\w+)"
	CNVMeDiskPrefix    = "nvme(\\d+)(n\\d+)?(p\\d+)?"
	CUSBDiskPattern    = "ub(\\w+)"
	CMMCDiskPattern    = "mmcblk(\\d+)(p\\d+)?"
	CVirtIODiskPattern = "vd(\\w+)"
	CXenDiskPattern    = "xvd(\\w+)"

	CFloppyDiskName = "Floppy"
	CCDROMDiskName  = "CD-ROM"
	CSCSIDiskName   = "SCSI"
	CIDEDiskName    = "IDE"
	CNVMeDiskName   = "NVMe"
	CUSBDiskName    = "USB"
	CMMCDiskName    = "MMC"
	CVirtIODiskName = "VirtIO"
	CXenDiskName    = "Xen"
)

var CFloppyDiskRegexp = regexp.MustCompile(CFloppyDiskPattern)
var CCDROMDiskRegexp = regexp.MustCompile(CCDROMDiskPattern)
var CSCSIDiskRegexp = regexp.MustCompile(CSCSIDiskPattern)
var CIDEDiskRegexp = regexp.MustCompile(CIDEDiskPattern)
var CNVMeDiskRegexp = regexp.MustCompile(CNVMeDiskPrefix)
var CUSBDiskRegexp = regexp.MustCompile(CUSBDiskPattern)
var CMMCDiskRegexp = regexp.MustCompile(CMMCDiskPattern)
var CVirtIODiskRegexp = regexp.MustCompile(CVirtIODiskPattern)
var CXenDiskRegexp = regexp.MustCompile(CXenDiskPattern)

var CDiskMap = map[string]regexp.Regexp{
	CFloppyDiskName: *CFloppyDiskRegexp,
	CCDROMDiskName:  *CCDROMDiskRegexp,
	CSCSIDiskName:   *CSCSIDiskRegexp,
	CIDEDiskName:    *CIDEDiskRegexp,
	CNVMeDiskName:   *CNVMeDiskRegexp,
	CUSBDiskName:    *CUSBDiskRegexp,
	CMMCDiskName:    *CMMCDiskRegexp,
	CVirtIODiskName: *CVirtIODiskRegexp,
	CXenDiskName:    *CXenDiskRegexp,
}
