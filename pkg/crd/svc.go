package crd

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	apiv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	clientv1alpha1 "github.com/onmetal/k8s-inventory/clientset/v1alpha1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/onmetal/inventory/pkg/inventory"
	"github.com/onmetal/inventory/pkg/lldp"
	"github.com/onmetal/inventory/pkg/netlink"
	"github.com/onmetal/inventory/pkg/utils"
)

const (
	CMACAddressLabelPrefix = "machine.onmetal.de/mac-address-"
	CSonicNamespace        = "switch.onmetal.de"
)

type Svc struct {
	client clientv1alpha1.InventoryInterface
}

func NewSvc(kubeconfig string, namespace string) (*Svc, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read kubeconfig from path %s", kubeconfig)
	}

	if err := apiv1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return nil, errors.Wrap(err, "unable to add registered types to client scheme")
	}

	clientset, err := clientv1alpha1.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build clientset from config")
	}

	client := clientset.Inventories(namespace)

	return &Svc{
		client: client,
	}, nil
}

func (s *Svc) BuildAndSave(inv *inventory.Inventory) error {
	cr, err := s.Build(inv)
	if err != nil {
		return errors.Wrap(err, "unable to build inventory resource manifest")
	}

	if err := s.Save(cr); err != nil {
		return errors.Wrap(err, "unable to save inventory resource")
	}

	return nil
}

func (s *Svc) Build(inv *inventory.Inventory) (*apiv1alpha1.Inventory, error) {
	cr := &apiv1alpha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1alpha1.InventorySpec{},
	}

	setters := []func(*apiv1alpha1.Inventory, *inventory.Inventory){
		s.setSystem,
		s.setIPMIs,
		s.setBlocks,
		s.setMemory,
		s.setCPUs,
		s.setNICs,
		s.setVirt,
		s.setHost,
		s.setDistro,
	}

	for _, setter := range setters {
		setter(cr, inv)
	}

	return cr, nil
}

func (s *Svc) Save(inv *apiv1alpha1.Inventory) error {
	_, err := s.client.Create(context.Background(), inv, metav1.CreateOptions{})
	if err == nil {
		return nil
	}
	if !apierrors.IsAlreadyExists(err) {
		return errors.Wrap(err, "unhandled error on creation")
	}

	existing, err := s.client.Get(context.Background(), inv.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to get resource")
	}

	existing.Spec = inv.Spec

	if existing.Labels == nil {
		existing.Labels = inv.Labels
	} else {
		for k, v := range inv.Labels {
			existing.Labels[k] = v
		}
	}

	if _, err := s.client.Update(context.Background(), existing, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(err, "unhandled error on update")
	}

	return nil
}

func (s *Svc) setSystem(cr *apiv1alpha1.Inventory, inv *inventory.Inventory) {
	if inv.DMI == nil {
		return
	}

	dmi := inv.DMI
	if dmi.SystemInformation == nil {
		return
	}

	if inv.Host == nil {
		return
	}

	// SONiC switches has dumb UUIDs like 03000200-0400-0500-0006-000700080009, maybe
	// the same on any switch, so it was decided to use md5 hash of serial number as UUID
	hostUUID := dmi.SystemInformation.UUID
	if inv.Host.Type == utils.CSwitchType {
		hostUUID = getUUID(CSonicNamespace, dmi.SystemInformation.SerialNumber)
	}
	cr.Name = hostUUID

	cr.Spec.System = &apiv1alpha1.SystemSpec{
		ID:           hostUUID,
		Manufacturer: dmi.SystemInformation.Manufacturer,
		ProductSKU:   dmi.SystemInformation.SKUNumber,
		SerialNumber: dmi.SystemInformation.SerialNumber,
	}
}

func (s *Svc) setIPMIs(cr *apiv1alpha1.Inventory, inv *inventory.Inventory) {
	ipmiDevCount := len(inv.IPMIDevices)
	if ipmiDevCount == 0 {
		return
	}

	ipmis := make([]apiv1alpha1.IPMISpec, ipmiDevCount)

	for i, ipmiDev := range inv.IPMIDevices {
		ipmi := apiv1alpha1.IPMISpec{
			IPAddress:  ipmiDev.IPAddress,
			MACAddress: ipmiDev.MACAddress,
		}

		ipmis[i] = ipmi
	}

	sort.Slice(ipmis, func(i, j int) bool {
		iStr := ipmis[i].MACAddress + ipmis[i].IPAddress
		jStr := ipmis[j].MACAddress + ipmis[j].IPAddress
		return iStr < jStr
	})

	cr.Spec.IPMIs = ipmis
}

func (s *Svc) setBlocks(cr *apiv1alpha1.Inventory, inv *inventory.Inventory) {
	if len(inv.BlockDevices) == 0 {
		return
	}

	blocks := make([]apiv1alpha1.BlockSpec, 0)
	var capacity uint64 = 0

	for _, blockDev := range inv.BlockDevices {
		// Filter non physical devices
		if blockDev.Type == "" {
			continue
		}

		var partitionTable *apiv1alpha1.PartitionTableSpec
		if blockDev.PartitionTable != nil {
			table := blockDev.PartitionTable

			partitionTable = &apiv1alpha1.PartitionTableSpec{
				Type: string(table.Type),
			}

			partCount := len(table.Partitions)
			if partCount > 0 {
				parts := make([]apiv1alpha1.PartitionSpec, partCount)

				for i, partition := range table.Partitions {
					part := apiv1alpha1.PartitionSpec{
						ID:   partition.ID,
						Name: partition.Name,
						Size: partition.Size,
					}

					parts[i] = part
				}

				sort.Slice(parts, func(i, j int) bool {
					return parts[i].ID < parts[j].ID
				})

				partitionTable.Partitions = parts
			}
		}

		block := apiv1alpha1.BlockSpec{
			Name:           blockDev.Name,
			Type:           blockDev.Type,
			Rotational:     blockDev.Rotational,
			Bus:            blockDev.Vendor,
			Model:          blockDev.Model,
			Size:           blockDev.Size,
			PartitionTable: partitionTable,
		}

		capacity += blockDev.Size
		blocks = append(blocks, block)
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Name < blocks[j].Name
	})

	cr.Spec.Blocks = &apiv1alpha1.BlockTotalSpec{
		Count:    uint64(len(blocks)),
		Capacity: capacity,
		Blocks:   blocks,
	}
}

func (s *Svc) setMemory(cr *apiv1alpha1.Inventory, inv *inventory.Inventory) {
	if inv.MemInfo == nil {
		return
	}

	cr.Spec.Memory = &apiv1alpha1.MemorySpec{
		Total: inv.MemInfo.MemTotal,
	}
}

func (s *Svc) setCPUs(cr *apiv1alpha1.Inventory, inv *inventory.Inventory) {
	if len(inv.CPUInfo) == 0 {
		return
	}

	cpuTotal := &apiv1alpha1.CPUTotalSpec{}

	cpuMarkMap := make(map[uint64]apiv1alpha1.CPUSpec, 0)

	for _, cpuInfo := range inv.CPUInfo {
		if val, ok := cpuMarkMap[cpuInfo.PhysicalID]; ok {
			val.LogicalIDs = append(val.LogicalIDs, cpuInfo.Processor)
			cpuMarkMap[cpuInfo.PhysicalID] = val
			continue
		}

		cpuTotal.Sockets += 1
		cpuTotal.Cores += cpuInfo.CpuCores
		cpuTotal.Threads += cpuInfo.Siblings

		cpu := apiv1alpha1.CPUSpec{
			PhysicalID:      cpuInfo.PhysicalID,
			LogicalIDs:      []uint64{cpuInfo.Processor},
			Cores:           cpuInfo.CpuCores,
			Siblings:        cpuInfo.Siblings,
			VendorID:        cpuInfo.VendorID,
			Family:          cpuInfo.CPUFamily,
			Model:           cpuInfo.Model,
			ModelName:       cpuInfo.ModelName,
			Stepping:        cpuInfo.Stepping,
			Microcode:       cpuInfo.Microcode,
			MHz:             resource.MustParse(cpuInfo.CPUMHz),
			CacheSize:       cpuInfo.CacheSize,
			FPU:             cpuInfo.FPU,
			FPUException:    cpuInfo.FPUException,
			CPUIDLevel:      cpuInfo.CPUIDLevel,
			WP:              cpuInfo.WP,
			Flags:           cpuInfo.Flags,
			VMXFlags:        cpuInfo.VMXFlags,
			Bugs:            cpuInfo.Bugs,
			BogoMIPS:        resource.MustParse(cpuInfo.BogoMIPS),
			CLFlushSize:     cpuInfo.CLFlushSize,
			CacheAlignment:  cpuInfo.CacheAlignment,
			AddressSizes:    cpuInfo.AddressSizes,
			PowerManagement: cpuInfo.PowerManagement,
		}
		sort.Strings(cpu.Flags)
		sort.Strings(cpu.VMXFlags)
		sort.Strings(cpu.Bugs)
		sort.Slice(cpu.LogicalIDs, func(i, j int) bool {
			return cpu.LogicalIDs[i] < cpu.LogicalIDs[j]
		})

		cpuMarkMap[cpuInfo.PhysicalID] = cpu
	}

	cpus := make([]apiv1alpha1.CPUSpec, 0)
	for _, v := range cpuMarkMap {
		cpus = append(cpus, v)
	}

	sort.Slice(cpus, func(i, j int) bool {
		return cpus[i].PhysicalID < cpus[j].PhysicalID
	})

	cpuTotal.CPUs = cpus

	cr.Spec.CPUs = cpuTotal
}

func (s *Svc) setNICs(cr *apiv1alpha1.Inventory, inv *inventory.Inventory) {
	if len(inv.NICs) == 0 {
		return
	}

	if inv.Host == nil {
		return
	}

	lldpMap := make(map[int][]apiv1alpha1.LLDPSpec)
	for _, frame := range inv.LLDPFrames {
		checkMap := make(map[lldp.Capability]struct{})
		enabledCapabilities := make([]apiv1alpha1.LLDPCapabilities, 0)
		for _, capability := range frame.EnabledCapabilities {
			if _, ok := checkMap[capability]; !ok {
				enabledCapabilities = append(enabledCapabilities, apiv1alpha1.LLDPCapabilities(capability))
				checkMap[capability] = struct{}{}
			}
		}
		sort.Slice(enabledCapabilities, func(i, j int) bool {
			return enabledCapabilities[i] < enabledCapabilities[j]
		})
		id, _ := strconv.Atoi(frame.InterfaceID)
		l := apiv1alpha1.LLDPSpec{
			ChassisID:         frame.ChassisID,
			SystemName:        frame.SystemName,
			SystemDescription: frame.SystemDescription,
			PortID:            frame.PortID,
			PortDescription:   frame.PortDescription,
			Capabilities:      enabledCapabilities,
		}

		if _, ok := lldpMap[id]; !ok {
			lldpMap[id] = make([]apiv1alpha1.LLDPSpec, 1)
		}

		lldpMap[id] = append(lldpMap[id], l)
	}

	ndpMap := make(map[int][]apiv1alpha1.NDPSpec)
	for _, ndp := range inv.NDPFrames {
		// filtering no arp as ip neigh does
		if ndp.State == netlink.CNeighbourNoARPCacheState {
			continue
		}

		n := apiv1alpha1.NDPSpec{
			IPAddress:  ndp.IP,
			MACAddress: ndp.MACAddress,
			State:      string(ndp.State),
		}

		if _, ok := ndpMap[ndp.DeviceIndex]; !ok {
			ndpMap[ndp.DeviceIndex] = make([]apiv1alpha1.NDPSpec, 1)
		}

		ndpMap[ndp.DeviceIndex] = append(ndpMap[ndp.DeviceIndex], n)
	}

	labels := cr.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	nics := make([]apiv1alpha1.NICSpec, 0)
	for _, nic := range inv.NICs {
		// filter non-physical interfaces according to type of inventorying host
		if inv.Host.Type == utils.CSwitchType {
			if nic.PCIAddress == "" && !strings.HasPrefix(nic.Name, "Ethernet") {
				continue
			}
		} else {
			if nic.PCIAddress == "" {
				continue
			}
		}

		lldps := lldpMap[int(nic.InterfaceIndex)]
		sort.Slice(lldps, func(i, j int) bool {
			iStr := lldps[i].ChassisID + lldps[i].PortID
			jStr := lldps[j].ChassisID + lldps[j].PortID
			return iStr < jStr
		})
		ndps := ndpMap[int(nic.InterfaceIndex)]
		sort.Slice(ndps, func(i, j int) bool {
			iStr := ndps[i].MACAddress + ndps[i].IPAddress
			jStr := ndps[j].MACAddress + ndps[j].IPAddress
			return iStr < jStr
		})

		ns := apiv1alpha1.NICSpec{
			Name:       nic.Name,
			PCIAddress: nic.PCIAddress,
			MACAddress: nic.Address,
			MTU:        nic.MTU,
			Speed:      nic.Speed,
			LLDPs:      lldps,
			NDPs:       ndps,
		}

		// Due to k8s validation which allows labels to consist of alphanumeric characters, '-', '_' or '.' need to replace
		// colons in nic's MAC address
		labels[CMACAddressLabelPrefix+strings.ToLower(nic.Name)] = strings.ReplaceAll(nic.Address, ":", "-")

		nics = append(nics, ns)
	}

	cr.SetLabels(labels)

	sort.Slice(nics, func(i, j int) bool {
		return nics[i].Name < nics[j].Name
	})

	cr.Spec.NICs = &apiv1alpha1.NICTotalSpec{
		Count: uint64(len(nics)),
		NICs:  nics,
	}
}

func (s *Svc) setVirt(cr *apiv1alpha1.Inventory, inv *inventory.Inventory) {
	if inv.Virtualization == nil {
		return
	}

	cr.Spec.Virt = &apiv1alpha1.VirtSpec{
		VMType: string(inv.Virtualization.Type),
	}
}

func (s *Svc) setHost(cr *apiv1alpha1.Inventory, inv *inventory.Inventory) {
	if inv.Host == nil {
		return
	}
	cr.Spec.Host = &apiv1alpha1.HostSpec{
		Type: inv.Host.Type,
		Name: inv.Host.Name,
	}
}

func (s *Svc) setDistro(cr *apiv1alpha1.Inventory, inv *inventory.Inventory) {
	if inv.Distro == nil {
		return
	}
	cr.Spec.Distro = &apiv1alpha1.DistroSpec{
		BuildVersion:  inv.Distro.BuildVersion,
		DebianVersion: inv.Distro.DebianVersion,
		KernelVersion: inv.Distro.KernelVersion,
		AsicType:      inv.Distro.AsicType,
		CommitId:      inv.Distro.CommitId,
		BuildDate:     inv.Distro.BuildDate,
		BuildNumber:   inv.Distro.BuildNumber,
		BuildBy:       inv.Distro.BuildBy,
	}
}

func getUUID(namespace string, identifier string) string {
	namespaceUUID := uuid.NewMD5(uuid.UUID{}, []byte(namespace))
	newUUID := uuid.NewMD5(namespaceUUID, []byte(identifier))
	return newUUID.String()
}
