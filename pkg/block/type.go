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
