// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"github.com/pkg/errors"
)

const (
	CMemInfoMemTotalKey          = "MemTotal"
	CMemInfoMemFreeKey           = "MemFree"
	CMemInfoMemAvailableKey      = "MemAvailable"
	CMemInfoBuffersKey           = "Buffers"
	CMemInfoCachedKey            = "Cached"
	CMemInfoSwapCachedKey        = "SwapCached"
	CMemInfoActiveKey            = "Active"
	CMemInfoInactiveKey          = "Inactive"
	CMemInfoActiveAnonKey        = "Active(anon)"
	CMemInfoInactiveAnonKey      = "Inactive(anon)"
	CMemInfoActiveFileKey        = "Active(file)"
	CMemInfoInactiveFileKey      = "Inactive(file)"
	CMemInfoUnevictableKey       = "Unevictable"
	CMemInfoMlockedKey           = "Mlocked"
	CMemInfoHighTotalKey         = "HighTotal"
	CMemInfoHighFreeKey          = "HighFree"
	CMemInfoLowTotalKey          = "LowTotal"
	CMemInfoLowFreeKey           = "LowFree"
	CMemInfoMmapCopyKey          = "MmapCopy"
	CMemInfoSwapTotalKey         = "SwapTotal"
	CMemInfoSwapFreeKey          = "SwapFree"
	CMemInfoDirtyKey             = "Dirty"
	CMemInfoWritebackKey         = "Writeback"
	CMemInfoAnonPagesKey         = "AnonPages"
	CMemInfoMappedKey            = "Mapped"
	CMemInfoShmemKey             = "Shmem"
	CMemInfoKReclaimableKey      = "KReclaimable"
	CMemInfoSlabKey              = "Slab"
	CMemInfoSReclaimableKey      = "SReclaimable"
	CMemInfoSUnreclaimKey        = "SUnreclaim"
	CMemInfoKernelStackKey       = "KernelStack"
	CMemInfoPageTablesKey        = "PageTables"
	CMemInfoQuicklistsKey        = "Quicklists"
	CMemInfoNFS_UnstableKey      = "NFS_Unstable"
	CMemInfoBounceKey            = "Bounce"
	CMemInfoWritebackTmpKey      = "WritebackTmp"
	CMemInfoCommitLimitKey       = "CommitLimit"
	CMemInfoCommitted_ASKey      = "Committed_AS"
	CMemInfoVmallocTotalKey      = "VmallocTotal"
	CMemInfoVmallocUsedKey       = "VmallocUsed"
	CMemInfoVmallocChunkKey      = "VmallocChunk"
	CMemInfoHardwareCorruptedKey = "HardwareCorrupted"
	CMemInfoLazyFreeKey          = "LazyFree"
	CMemInfoAnonHugePagesKey     = "AnonHugePages"
	CMemInfoShmemHugePagesKey    = "ShmemHugePages"
	CMemInfoShmemPmdMappedKey    = "ShmemPmdMapped"
	CMemInfoCmaTotalKey          = "CmaTotal"
	CMemInfoCmaFreeKey           = "CmaFree"
	CMemInfoHugePages_TotalKey   = "HugePages_Total"
	CMemInfoHugePages_FreeKey    = "HugePages_Free"
	CMemInfoHugePages_RsvdKey    = "HugePages_Rsvd"
	CMemInfoHugePages_SurpKey    = "HugePages_Surp"
	CMemInfoHugepagesizeKey      = "Hugepagesize"
	CMemInfoDirectMap4kKey       = "DirectMap4k"
	CMemInfoDirectMap4MKey       = "DirectMap4M"
	CMemInfoDirectMap2MKey       = "DirectMap2M"
	CMemInfoDirectMap1GKey       = "DirectMap1G"
	// Undocumented, but in kernel
	CMemInfoPercpuKey        = "Percpu"
	CMemInfoFileHugePagesKey = "FileHugePages"
	CMemInfoFilePmdMappedKey = "FilePmdMapped"
	CMemInfoHugetlbKey       = "Hugetlb"
	// NUMA specific
	CMemInfoMemUsedKey   = "MemUsed"
	CMemInfoFilePagesKey = "FilePages"
)

type Info struct {
	MemTotal          uint64
	MemFree           uint64
	MemAvailable      uint64
	Buffers           uint64
	Cached            uint64
	SwapCached        uint64
	Active            uint64
	Inactive          uint64
	ActiveAnon        uint64
	InactiveAnon      uint64
	ActiveFile        uint64
	InactiveFile      uint64
	Unevictable       uint64
	Mlocked           uint64
	HighTotal         uint64
	HighFree          uint64
	LowTotal          uint64
	LowFree           uint64
	MmapCopy          uint64
	SwapTotal         uint64
	SwapFree          uint64
	Dirty             uint64
	Writeback         uint64
	AnonPages         uint64
	Mapped            uint64
	Shmem             uint64
	KReclaimable      uint64
	Slab              uint64
	SReclaimable      uint64
	SUnreclaim        uint64
	KernelStack       uint64
	PageTables        uint64
	Quicklists        uint64
	NFS_Unstable      uint64
	Bounce            uint64
	WritebackTmp      uint64
	CommitLimit       uint64
	Committed_AS      uint64
	VmallocTotal      uint64
	VmallocUsed       uint64
	VmallocChunk      uint64
	HardwareCorrupted uint64
	LazyFree          uint64
	AnonHugePages     uint64
	ShmemHugePages    uint64
	ShmemPmdMapped    uint64
	CmaTotal          uint64
	CmaFree           uint64
	HugePages_Total   uint64
	HugePages_Free    uint64
	HugePages_Rsvd    uint64
	HugePages_Surp    uint64
	Hugepagesize      uint64
	DirectMap4k       uint64
	DirectMap4M       uint64
	DirectMap2M       uint64
	DirectMap1G       uint64
	// Undocumented, but in kernel
	Percpu        uint64
	FileHugePages uint64
	FilePmdMapped uint64
	Hugetlb       uint64
	// NUMA specific
	MemUsed   uint64
	FilePages uint64
}

func (mem *Info) setField(key string, val uint64) error {
	switch key {
	case CMemInfoMemTotalKey:
		mem.MemTotal = val
	case CMemInfoMemFreeKey:
		mem.MemFree = val
	case CMemInfoMemAvailableKey:
		mem.MemAvailable = val
	case CMemInfoBuffersKey:
		mem.Buffers = val
	case CMemInfoCachedKey:
		mem.Cached = val
	case CMemInfoSwapCachedKey:
		mem.SwapCached = val
	case CMemInfoActiveKey:
		mem.Active = val
	case CMemInfoInactiveKey:
		mem.Inactive = val
	case CMemInfoActiveAnonKey:
		mem.ActiveAnon = val
	case CMemInfoInactiveAnonKey:
		mem.InactiveAnon = val
	case CMemInfoActiveFileKey:
		mem.ActiveFile = val
	case CMemInfoInactiveFileKey:
		mem.InactiveFile = val
	case CMemInfoUnevictableKey:
		mem.Unevictable = val
	case CMemInfoMlockedKey:
		mem.Mlocked = val
	case CMemInfoHighTotalKey:
		mem.HighTotal = val
	case CMemInfoHighFreeKey:
		mem.HighFree = val
	case CMemInfoLowTotalKey:
		mem.LowTotal = val
	case CMemInfoLowFreeKey:
		mem.LowFree = val
	case CMemInfoMmapCopyKey:
		mem.MmapCopy = val
	case CMemInfoSwapTotalKey:
		mem.SwapTotal = val
	case CMemInfoSwapFreeKey:
		mem.SwapFree = val
	case CMemInfoDirtyKey:
		mem.Dirty = val
	case CMemInfoWritebackKey:
		mem.Writeback = val
	case CMemInfoAnonPagesKey:
		mem.AnonPages = val
	case CMemInfoMappedKey:
		mem.Mapped = val
	case CMemInfoShmemKey:
		mem.Shmem = val
	case CMemInfoKReclaimableKey:
		mem.KReclaimable = val
	case CMemInfoSlabKey:
		mem.Slab = val
	case CMemInfoSReclaimableKey:
		mem.SReclaimable = val
	case CMemInfoSUnreclaimKey:
		mem.SUnreclaim = val
	case CMemInfoKernelStackKey:
		mem.KernelStack = val
	case CMemInfoPageTablesKey:
		mem.PageTables = val
	case CMemInfoQuicklistsKey:
		mem.Quicklists = val
	case CMemInfoNFS_UnstableKey:
		mem.NFS_Unstable = val
	case CMemInfoBounceKey:
		mem.Bounce = val
	case CMemInfoWritebackTmpKey:
		mem.WritebackTmp = val
	case CMemInfoCommitLimitKey:
		mem.CommitLimit = val
	case CMemInfoCommitted_ASKey:
		mem.Committed_AS = val
	case CMemInfoVmallocTotalKey:
		mem.VmallocTotal = val
	case CMemInfoVmallocUsedKey:
		mem.VmallocUsed = val
	case CMemInfoVmallocChunkKey:
		mem.VmallocChunk = val
	case CMemInfoHardwareCorruptedKey:
		mem.HardwareCorrupted = val
	case CMemInfoLazyFreeKey:
		mem.LazyFree = val
	case CMemInfoAnonHugePagesKey:
		mem.AnonHugePages = val
	case CMemInfoShmemHugePagesKey:
		mem.ShmemHugePages = val
	case CMemInfoShmemPmdMappedKey:
		mem.ShmemPmdMapped = val
	case CMemInfoCmaTotalKey:
		mem.CmaTotal = val
	case CMemInfoCmaFreeKey:
		mem.CmaFree = val
	case CMemInfoHugePages_TotalKey:
		mem.HugePages_Total = val
	case CMemInfoHugePages_FreeKey:
		mem.HugePages_Free = val
	case CMemInfoHugePages_RsvdKey:
		mem.HugePages_Rsvd = val
	case CMemInfoHugePages_SurpKey:
		mem.HugePages_Surp = val
	case CMemInfoHugepagesizeKey:
		mem.Hugepagesize = val
	case CMemInfoDirectMap4kKey:
		mem.DirectMap4k = val
	case CMemInfoDirectMap4MKey:
		mem.DirectMap4M = val
	case CMemInfoDirectMap2MKey:
		mem.DirectMap2M = val
	case CMemInfoDirectMap1GKey:
		mem.DirectMap1G = val
	// Undocumented, but in kernel
	case CMemInfoPercpuKey:
		mem.Percpu = val
	case CMemInfoFileHugePagesKey:
		mem.FileHugePages = val
	case CMemInfoFilePmdMappedKey:
		mem.FilePmdMapped = val
	case CMemInfoHugetlbKey:
		mem.Hugetlb = val
	// NUMA specific
	case CMemInfoMemUsedKey:
		mem.MemUsed = val
	case CMemInfoFilePagesKey:
		mem.FilePages = val
	default:
		return errors.Errorf("unknown key %s from meminfo", key)
	}
	return nil
}
