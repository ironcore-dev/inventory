package lldp

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/file"
	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/utils"
)

const (
	CLLDPPath     = "/run/systemd/netif/lldp"
	CClassNetPath = "/sys/class/net/"
	CIndexFile    = "ifindex"
)

type Svc struct {
	printer      *printer.Svc
	frameInfoSvc *FrameSvc
	hostSvc      *host.Svc
	lldpPath     string
	indexPath    string
}

func NewSvc(printer *printer.Svc, frameInfoSvc *FrameSvc, hostSvc *host.Svc, basePath string) *Svc {
	return &Svc{
		printer:      printer,
		frameInfoSvc: frameInfoSvc,
		hostSvc:      hostSvc,
		lldpPath:     path.Join(basePath, CLLDPPath),
		indexPath:    path.Join(basePath, CClassNetPath),
	}
}

func (s *Svc) GetData() ([]Frame, error) {
	frameInfos := make([]Frame, 0)

	hostInfo, err := s.hostSvc.GetData()
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "failed to collect host info"))
	}

	switch hostInfo.Type {
	case utils.CMachineType:
		frameFiles, err := ioutil.ReadDir(s.lldpPath)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get list of frame files")
		}
		// iterate over /run/systemd/netif/lldp/%i
		for _, frameFile := range frameFiles {
			fName := frameFile.Name()
			filePath := path.Join(s.lldpPath, fName)
			info, err := s.frameInfoSvc.GetFrame(fName, filePath)
			if err != nil {
				s.printer.VErr(errors.Errorf("unable to collect LLDP info for interface idx %s", fName))
				continue
			}
			frameInfos = append(frameInfos, *info)
		}
	case utils.CSwitchType:
		redisClient := utils.Client()
		lldpKeys, err := utils.GetKeysByPattern(redisClient, utils.CLLDPEntryKeyMask)
		if err != nil {
			s.printer.VErr(err)
			return nil, err
		}
		for _, key := range lldpKeys {
			frame, err := s.processRedisPortData(key, redisClient)
			if err != nil {
				return nil, err
			}
			frameInfos = append(frameInfos, *frame)
		}
	}
	return frameInfos, nil
}

func (s *Svc) processRedisPortData(key string, redisClient *redis.Client) (*Frame, error) {
	port := strings.Split(key, ":")
	filePath := path.Join(s.indexPath, port[1], CIndexFile)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get interface index value from %s", filePath)
	}
	rawData, err := utils.GetValuesFromHashEntry(redisClient, key, &utils.CRedisLLDPFields)
	if err != nil {
		s.printer.VErr(errors.Errorf("unable to collect LLDP info for interface %s", port[1]))
		return nil, err
	}
	capabilities, err := GetCapabilities(rawData["lldp_rem_sys_cap_supported"])
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to collect supported capabilities for remote interface"))
		return nil, err
	}
	enabledCapabilities, err := GetCapabilities(rawData["lldp_rem_sys_cap_enabled"])
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to collect enabled capabilities for remote interface"))
		return nil, err
	}
	frame := &Frame{
		InterfaceID:         fileVal,
		ChassisID:           rawData["lldp_rem_chassis_id"],
		SystemName:          rawData["lldp_rem_sys_name"],
		SystemDescription:   rawData["lldp_rem_sys_desc"],
		Capabilities:        capabilities,
		EnabledCapabilities: enabledCapabilities,
		PortID:              rawData["lldp_rem_port_id"],
		PortDescription:     rawData["lldp_rem_port_desc"],
		ManagementAddresses: strings.Split(rawData["lldp_rem_man_addr"], ","),
		TTL:                 0,
	}
	return frame, nil
}
