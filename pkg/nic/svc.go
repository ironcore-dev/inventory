package nic

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/prometheus/procfs/sysfs"

	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/redis"
	"github.com/onmetal/inventory/pkg/utils"

	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
)

const (
	CNICDevicePath = "/sys/class/net"
)

type Svc struct {
	printer    *printer.Svc
	nicDevSvc  *DeviceSvc
	nicDevPath string
	hostSvc    *host.Svc
	redisSvc   *redis.Svc
}

func NewSvc(printer *printer.Svc, nicDevSvc *DeviceSvc, hostSvc *host.Svc, redisSvc *redis.Svc, basePath string) *Svc {
	return &Svc{
		printer:    printer,
		nicDevSvc:  nicDevSvc,
		hostSvc:    hostSvc,
		redisSvc:   redisSvc,
		nicDevPath: basePath,
	}
}

func (s *Svc) GetData() ([]Device, error) {
	hostInfo, err := s.hostSvc.GetData()
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "failed to collect host info"))
	}

	nicFolders, err := os.ReadDir(path.Join(s.nicDevPath, CNICDevicePath))
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of nic folders")
	}

	fs, err := sysfs.NewFS(path.Join(s.nicDevPath, sysfs.DefaultMountPoint))
	if err != nil {
		return nil, errors.Wrap(err, "unable to init sysfs")
	}

	var nics []Device
	for _, nicFolder := range nicFolders {
		if nicFolder.Type().IsRegular() {
			continue
		}

		fName := nicFolder.Name()
		netClassIface, err := fs.NetClassByIface(fName)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to collect Device data"))
			continue
		}

		nic := &Device{
			Name:                netClassIface.Name,
			Address:             netClassIface.Address,
			AddressLength:       *netClassIface.AddrLen,
			Broadcast:           netClassIface.Broadcast,
			Carrier:             false,
			CarrierChanges:      *netClassIface.CarrierChanges,
			CarrierDownCount:    *netClassIface.CarrierDownCount,
			CarrierUpCount:      *netClassIface.CarrierUpCount,
			DevID:               fmt.Sprintf("%d", netClassIface.DevID),
			Dormant:             false,
			Duplex:              netClassIface.Duplex,
			InterfaceAlias:      netClassIface.IfAlias,
			InterfaceIndex:      *netClassIface.IfIndex,
			InterfaceLink:       *netClassIface.IfLink,
			MTU:                 uint16(*netClassIface.MTU),
			NetDevGroup:         *netClassIface.NetDevGroup,
			OperationalState:    netClassIface.OperState,
			PhysicalPortID:      netClassIface.PhysPortID,
			PhysicalPortName:    netClassIface.PhysPortName,
			PhysicalSwitchID:    netClassIface.PhysSwitchID,
			TransmitQueueLength: *netClassIface.TxQueueLen,
		}

		err = nic.defPCIAddress(path.Join(s.nicDevPath, CNICDevicePath, fName))
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to get pci address"))
			continue
		}

		if netClassIface.AddrAssignType != nil && *netClassIface.AddrAssignType >= int64(len(CAddressAssignTypes)) {
			nic.AddressAssignType = CAddressAssignTypes[*netClassIface.AddrAssignType]
		}

		if netClassIface.NameAssignType != nil && *netClassIface.NameAssignType >= int64(len(CNameAssignTypes)) {
			nic.NameAssignType = CNameAssignTypes[*netClassIface.NameAssignType]
		}

		if netClassIface.LinkMode != nil && *netClassIface.LinkMode >= int64(len(CLinkModes)) {
			nic.LinkMode = CLinkModes[*netClassIface.LinkMode]
		}

		if netClassIface.Carrier != nil && *netClassIface.Carrier > 0 {
			nic.Carrier = true
		}

		if netClassIface.Dormant != nil && *netClassIface.Dormant > 0 {
			nic.Dormant = true
		}

		nic.Flags = NewFlags(*netClassIface.Flags)

		if netClassIface.Type != nil && *netClassIface.Type >= int64(len(CTypes)) {
			theType, ok := CTypes[*netClassIface.Type]
			if ok {
				nic.Type = theType
			} else {
				nic.Type = CTypes[0xffff]
			}
		} else {
			nic.Type = CTypes[0xffff]
		}

		if hostInfo.Type == utils.CSwitchType && strings.HasPrefix(fName, "Ethernet") {
			info, err := s.redisSvc.GetPortAdditionalInfo(fName)
			if err != nil {
				s.printer.VErr(errors.Wrap(err, "unable to collect additional Device data from Redis"))
				continue
			}
			nic.Lanes = uint8(len(strings.Split(info[redis.CPortLanes], ",")))
			nic.FEC = info[redis.CPortFec]
			if nic.FEC == "" {
				nic.FEC = switchv1beta1.CFECNone
			}
			speed, err := strconv.Atoi(info[redis.CPortSpeed])
			if err != nil {
				s.printer.VErr(errors.Wrap(err, "unable to collect additional Device data from Redis"))
				continue
			}
			nic.Speed = uint32(speed)
		}
		nics = append(nics, *nic)
	}

	return nics, nil
}
