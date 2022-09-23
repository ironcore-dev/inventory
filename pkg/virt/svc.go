package virt

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/jeek120/cpuid"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/cpu"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/file"
)

const (
	CTypeNone        = "none"
	CTypeVMKVM       = "kvm"
	CTypeVMQEMU      = "qemu"
	CTypeVMBochs     = "bochs"
	CTypeVMXen       = "xen"
	CTypeVMUML       = "uml"
	CTypeVMVMware    = "vmware"
	CTypeVMOracle    = "oracle"
	CTypeVMMicrosoft = "microsoft"
	CTypeVMZVM       = "zvm"
	CTypeVMParallels = "parallels"
	CTypeVMBhyve     = "bhyve"
	CTypeVMQNX       = "qnx"
	CTypeVMACRN      = "acrn"
	CTypeVMPowerVM   = "powervm"
	CTypeVMOther     = "other"

	CProcXenPath                        = "/proc/xen"
	CSysHypervisorTypePath              = "/sys/hypervisor/type"
	CDeviceTreePath                     = "/proc/device-tree"
	CDeviceTreeHypervisorCompatiblePath = "/proc/device-tree/hypervisor/compatible"
	CDeviceTreeIBMPartitionNamePath     = "/proc/device-tree/ibm,partition-name"
	CDeviceTreeHMCManagedPath           = "/proc/device-tree/hmc-managed?"
	CDeviceTreeQEMUPath                 = "/proc/device-tree/chosen/qemu,graphic-width"
	CProcSysInfoPath                    = "/proc/sysinfo"
	CSysHypervisorFeaturesPath          = "/sys/hypervisor/properties/features"
	CProcXenCapabilitiesPath            = "/proc/xen/capabilities"
)

type Type string

type Virtualization struct {
	Type Type
}

type Svc struct {
	dmiSvc     *dmi.Svc
	cpuInfoSvc *cpu.InfoSvc

	procXenPath                        string
	sysHypervisorTypePath              string
	deviceTreePath                     string
	deviceTreeHypervisorCompatiblePath string
	deviceTreeIBMPartitionNamePath     string
	deviceTreeHMCManagedPath           string
	deviceTreeQEMUPath                 string
	procSysInfoPath                    string
	sysHypervisorFeaturesPath          string
	procXenCapabilitiesPath            string
}

func NewSvc(dmiSvc *dmi.Svc, cpuInfoSvc *cpu.InfoSvc, basePath string) *Svc {
	return &Svc{
		dmiSvc:                             dmiSvc,
		cpuInfoSvc:                         cpuInfoSvc,
		procXenPath:                        path.Join(basePath, CProcXenPath),
		sysHypervisorTypePath:              path.Join(basePath, CSysHypervisorTypePath),
		deviceTreePath:                     path.Join(basePath, CDeviceTreePath),
		deviceTreeHypervisorCompatiblePath: path.Join(basePath, CDeviceTreeHypervisorCompatiblePath),
		deviceTreeIBMPartitionNamePath:     path.Join(basePath, CDeviceTreeIBMPartitionNamePath),
		deviceTreeHMCManagedPath:           path.Join(basePath, CDeviceTreeHMCManagedPath),
		deviceTreeQEMUPath:                 path.Join(basePath, CDeviceTreeQEMUPath),
		procSysInfoPath:                    path.Join(basePath, CProcSysInfoPath),
		sysHypervisorFeaturesPath:          path.Join(basePath, CSysHypervisorFeaturesPath),
		procXenCapabilitiesPath:            path.Join(basePath, CProcXenCapabilitiesPath),
	}
}

var CDMIVendorPrefixes = map[string]Type{
	"KVM":                CTypeVMKVM,
	"QEMU":               CTypeVMQEMU,
	"VMware":             CTypeVMVMware,
	"VMW":                CTypeVMVMware,
	"innotek GmbH":       CTypeVMOracle,
	"Oracle Corporation": CTypeVMOracle,
	"Xen":                CTypeVMXen,
	"Bochs":              CTypeVMBochs,
	"Parallels":          CTypeVMParallels,
	"BHYVE":              CTypeVMBhyve,
}

var CCPUIDVMStrings = map[string]Type{
	"XenVMMXenVMM": CTypeVMXen,
	"KVMKVMKVM":    CTypeVMKVM,
	"TCGTCGTCGTCG": CTypeVMQEMU,
	"VMwareVMware": CTypeVMVMware,
	"Microsoft Hv": CTypeVMMicrosoft,
	"bhyve bhyve ": CTypeVMBhyve,
	"QNXQVMBSQG":   CTypeVMQNX,
	"ACRNACRNACRN": CTypeVMACRN,
}

func (s *Svc) GetData() (*Virtualization, error) {
	theType, err := s.getType()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to complete virt detection")
	}
	return &Virtualization{
		Type: theType,
	}, nil
}

func (s *Svc) getType() (Type, error) {
	// check dmi
	//   if oracle
	//     return result
	//	 if xen
	//	   check xen dom0
	//     return result
	// check cpuInfo
	//   if xen
	//     check xen dom0
	//     return result
	//   if !none && !other
	//     return result
	// check cpuId
	//   if xen
	//     check xen dom0
	//     return result
	//   if !none && !other
	//     return result
	// check dmi
	//   if !none && !other
	//     return result
	// check xen
	//	 if xen
	//	   check xen dom0
	//     return result
	//   if !none && !other
	//     return result
	// check hypervisor
	//	 if xen
	//	   check xen dom0
	//     return result
	//   if !none && !other
	//     return result
	// check devTree
	//	 if xen
	//	   check xen dom0
	//     return result
	//   if !none && !other
	//     return result
	// check zvm
	//	 if xen
	//	   check xen dom0
	//     return result
	//   return result

	dmiType, err := s.checkVMWithDMI()
	if err != nil {
		return "", errors.Wrapf(err, "unable to check for vm with dmi")
	}
	if dmiType == CTypeVMOracle {
		return dmiType, nil
	}
	if dmiType == CTypeVMXen {
		if theType, err := s.checkVMWithXenDom0(); err != nil {
			return "", errors.Wrapf(err, "unable to check for xen dom0")
		} else {
			return theType, nil
		}
	}

	passType := s.checkVMWithCPUInfo()
	if passType == CTypeVMXen {
		if theType, err := s.checkVMWithXenDom0(); err != nil {
			return "", errors.Wrapf(err, "unable to check for xen dom0")
		} else {
			return theType, nil
		}
	}
	if passType != CTypeNone && passType != CTypeVMOther {
		return passType, nil
	}

	passType, err = s.checkVMWithCPUID()
	if err != nil {
		return "", errors.Wrapf(err, "unable to check for cpuid")
	}
	if passType == CTypeVMXen {
		if theType, err := s.checkVMWithXenDom0(); err != nil {
			return "", errors.Wrapf(err, "unable to check for xen dom0")
		} else {
			return theType, nil
		}
	}
	if passType != CTypeNone && passType != CTypeVMOther {
		return passType, nil
	}

	if dmiType != CTypeNone && dmiType != CTypeVMOther {
		return dmiType, nil
	}

	passType = s.checkVMWithXen()
	if passType == CTypeVMXen {
		if theType, err := s.checkVMWithXenDom0(); err != nil {
			return "", errors.Wrapf(err, "unable to check for xen dom0")
		} else {
			return theType, nil
		}
	}
	if passType != CTypeNone && passType != CTypeVMOther {
		return passType, nil
	}

	passType, err = s.checkVMWithSysHypervisor()
	if err != nil {
		return "", errors.Wrapf(err, "unable to check for hypervizor")
	}
	if passType == CTypeVMXen {
		if theType, err := s.checkVMWithXenDom0(); err != nil {
			return "", errors.Wrapf(err, "unable to check for xen dom0")
		} else {
			return theType, nil
		}
	}
	if passType != CTypeNone && passType != CTypeVMOther {
		return passType, nil
	}

	passType, err = s.checkVMWithDeviceTree()
	if err != nil {
		return "", errors.Wrapf(err, "unable to check for device tree")
	}
	if passType == CTypeVMXen {
		if theType, err := s.checkVMWithXenDom0(); err != nil {
			return "", errors.Wrapf(err, "unable to check for xen dom0")
		} else {
			return theType, nil
		}
	}
	if passType != CTypeNone && passType != CTypeVMOther {
		return passType, nil
	}

	passType, err = s.checkVMWithZVM()
	if err != nil {
		return "", errors.Wrapf(err, "unable to check for zvm")
	}
	if passType == CTypeVMXen {
		if theType, err := s.checkVMWithXenDom0(); err != nil {
			return "", errors.Wrapf(err, "unable to check for xen dom0")
		} else {
			return theType, nil
		}
	}
	if passType != CTypeNone && passType != CTypeVMOther {
		return passType, nil
	}

	return passType, nil
}

func (s *Svc) checkVMWithDMI() (Type, error) {
	dmiData, err := s.dmiSvc.GetData()
	if dmiData == nil {
		return CTypeNone, errors.Wrapf(err, "unable to get valid dmiData")
	}

	vendorLocators := []string{}
	if dmiData.SystemInformation != nil {
		vendorLocators = append(vendorLocators, dmiData.SystemInformation.ProductName, dmiData.SystemInformation.Manufacturer)
	}
	if dmiData.BIOSInformation != nil {
		vendorLocators = append(vendorLocators, dmiData.BIOSInformation.Vendor)
	}
	for _, board := range dmiData.BoardInformation {
		vendorLocators = append(vendorLocators, board.Manufacturer)
	}

	for _, vendorLocator := range vendorLocators {
		for k, v := range CDMIVendorPrefixes {
			if strings.HasPrefix(vendorLocator, k) {
				return v, nil
			}
		}
	}

	return CTypeNone, nil
}

func (s *Svc) checkVMWithCPUInfo() Type {
	infos, _ := s.cpuInfoSvc.GetInfo()

	for _, info := range infos {
		if strings.HasPrefix(info.VendorID, "User Mode Linux") {
			return CTypeVMUML
		}
	}

	return CTypeNone
}

func (s *Svc) checkVMWithCPUID() (Type, error) {
	ids := [4]uint32{}
	cpuid.Cpuid(&ids, 1)

	hypervisor := ids[2] & 0x80000000

	if hypervisor == 0 {
		return CTypeNone, nil
	}

	ids = [4]uint32{}
	cpuid.Cpuid(&ids, 0x40000000)

	buf := new(bytes.Buffer)
	for i := 1; i < 4; i++ {
		if err := binary.Write(buf, binary.LittleEndian, ids[i]); err != nil {
			return "", errors.Wrap(err, "unable to write to buffer")
		}
	}

	text := buf.String()

	for k, v := range CCPUIDVMStrings {
		if strings.HasPrefix(text, k) {
			return v, nil
		}
	}

	return CTypeVMOther, nil
}

func (s *Svc) checkVMWithXen() Type {
	if _, err := os.Stat(s.procXenPath); os.IsNotExist(err) {
		return CTypeNone
	}

	return CTypeVMXen
}

func (s *Svc) checkVMWithSysHypervisor() (Type, error) {
	if _, err := os.Stat(s.sysHypervisorTypePath); os.IsNotExist(err) {
		return CTypeNone, nil
	}

	str, err := file.ToString(s.sysHypervisorTypePath)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read %s into string", s.sysHypervisorTypePath)
	}

	if str == "xen" {
		return CTypeVMXen, nil
	}

	return CTypeVMOther, nil
}

func (s *Svc) checkVMWithDeviceTree() (Type, error) {
	_, err := os.Stat(s.deviceTreeHypervisorCompatiblePath)
	if !os.IsNotExist(err) {
		return "", errors.Wrapf(err, "unable to open file %s", s.deviceTreeHypervisorCompatiblePath)
	} else {
		access := 0

		if _, err := os.Stat(s.deviceTreeIBMPartitionNamePath); err != nil {
			access++
		}
		if _, err := os.Stat(s.deviceTreeHMCManagedPath); err != nil {
			access++
		}
		if _, err := os.Stat(s.deviceTreeQEMUPath); err != nil {
			access++
		}

		if access == 0 {
			return CTypeVMPowerVM, nil
		}

		files, err := os.ReadDir(s.deviceTreePath)
		if os.IsNotExist(err) {
			return CTypeNone, nil
		}
		if err != nil {
			return "", errors.Wrapf(err, "unable to open directory %s", s.deviceTreePath)
		}

		for _, f := range files {
			if strings.Contains(f.Name(), "fw-cfg") {
				return CTypeVMQEMU, nil
			}
		}
	}

	str, err := file.ToString(s.deviceTreeHypervisorCompatiblePath)
	if err != nil {
		return "", errors.Wrapf(err, "unable to open file %s", s.deviceTreeHypervisorCompatiblePath)
	}

	if str == "linux,kvm" {
		return CTypeVMKVM, nil
	}
	if strings.Contains(str, "xen") {
		return CTypeVMXen, nil
	}
	if strings.Contains(str, "vmware") {
		return CTypeVMXen, nil
	}

	return CTypeVMOther, nil
}

func (s *Svc) checkVMWithZVM() (Type, error) {
	sysInfoData, err := os.ReadFile(s.procSysInfoPath)
	if os.IsNotExist(err) {
		return CTypeNone, nil
	}
	if err != nil {
		return "", errors.Wrapf(err, "unable to read cpu info from %s", s.procSysInfoPath)
	}

	bufReader := bytes.NewReader(sysInfoData)
	scanner := bufio.NewScanner(bufReader)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "VM00 Control Program") {
			continue
		}
		if strings.Contains(line, "z/VM") {
			return CTypeVMZVM, nil
		} else {
			return CTypeVMKVM, nil
		}
	}

	return CTypeNone, nil
}

func (s *Svc) checkVMWithXenDom0() (Type, error) {
	str, err := file.ToString(s.sysHypervisorFeaturesPath)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read %s", s.sysHypervisorFeaturesPath)
	}

	features, err := strconv.ParseUint("0x"+str, 16, 64)
	if err == nil {
		t := features & (1 << 11)
		if t == 0 {
			return CTypeVMXen, nil
		} else {
			return CTypeNone, nil
		}
	}

	str, err = file.ToString(s.sysHypervisorFeaturesPath)
	cause := errors.Unwrap(err)
	if !os.IsNotExist(cause) {
		return CTypeVMXen, nil
	}
	if err != nil {
		return "", errors.Wrapf(err, "unable to read %s", s.sysHypervisorFeaturesPath)
	}

	capabilities := strings.Split(str, ",")

	for _, capability := range capabilities {
		if capability == "control_d" {
			return CTypeNone, nil
		}
	}

	return CTypeVMXen, nil
}
