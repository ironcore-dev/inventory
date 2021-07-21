package redis

import (
	"context"
	"math/bits"
	"path"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/file"
	"github.com/onmetal/inventory/pkg/lldp/frame"
)

const (
	CLLDPEntryKeyMask = "LLDP_ENTRY*"
	CPortEntryPrefix  = "PORT_TABLE:"
	CClassNetPath     = "/sys/class/net/"
	CIndexFile        = "ifindex"
)

const (
	CLLDPRemoteChassisId             = "lldp_rem_chassis_id"
	CLLDPRemoteSystemName            = "lldp_rem_sys_name"
	CLLDPRemoteSystemDescription     = "lldp_rem_sys_desc"
	CLLDPRemoteCapabilitiesSupported = "lldp_rem_sys_cap_supported"
	CLLDPRemoteCapabilitiesEnabled   = "lldp_rem_sys_cap_enabled"
	CLLDPRemotePortId                = "lldp_rem_port_id"
	CLLDPRemotePortDescription       = "lldp_rem_port_desc"
	CLLDPRemoteManagementAddresses   = "lldp_rem_man_addr"
)

const (
	CPortLanes = "lanes"
	CPortFec   = "fec"
	CPortSpeed = "speed"
)

var CRedisLLDPFields = []string{
	CLLDPRemoteChassisId,
	CLLDPRemoteSystemName,
	CLLDPRemoteSystemDescription,
	CLLDPRemoteCapabilitiesSupported,
	CLLDPRemoteCapabilitiesEnabled,
	CLLDPRemotePortId,
	CLLDPRemotePortDescription,
	CLLDPRemoteManagementAddresses,
}

var CRedisPortFields = []string{
	CPortLanes,
	CPortFec,
	CPortSpeed,
}

type Svc struct {
	client    *redis.Client
	ctx       context.Context
	indexPath string
}

func NewRedisSvc(basePath string) *Svc {
	return &Svc{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		ctx:       context.Background(),
		indexPath: path.Join(basePath, CClassNetPath),
	}
}

func (s *Svc) GetFrames() ([]frame.Frame, error) {
	frames := make([]frame.Frame, 0)
	lldpKeys, err := s.getKeysByPattern(CLLDPEntryKeyMask)
	if err != nil {
		return nil, err
	}
	for _, key := range lldpKeys {
		f, err := s.processRedisPortData(key)
		if err != nil {
			return nil, err
		}
		frames = append(frames, *f)
	}
	return frames, nil
}

func (s *Svc) GetPortAdditionalInfo(name string) (map[string]string, error) {
	result := map[string]string{CPortLanes: "", CPortFec: ""}
	key := CPortEntryPrefix + name
	for _, f := range CRedisPortFields {
		val, err := s.client.Do(s.ctx, "HGET", key, f).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return nil, errors.Wrap(err, "failed to get value")
		}
		result[f] = val.(string)
	}
	return result, nil
}

func (s *Svc) getKeysByPattern(pattern string) ([]string, error) {
	val, err := s.client.Keys(s.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (s *Svc) getValuesFromHashEntry(key string, fields *[]string) (map[string]string, error) {
	result := make(map[string]string)
	for _, f := range *fields {
		val, err := s.client.Do(s.ctx, "HGET", key, f).Result()
		if err != nil {
			if err == redis.Nil {
				cause := errors.New("key not found")
				return nil, errors.Wrap(cause, key)
			}
			return nil, errors.Wrap(err, "failed to get value")
		}
		result[f] = val.(string)
	}
	return result, nil
}

func (s *Svc) processRedisPortData(key string) (*frame.Frame, error) {
	port := strings.Split(key, ":")
	filePath := path.Join(s.indexPath, port[1], CIndexFile)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get interface index value from %s", filePath)
	}
	rawData, err := s.getValuesFromHashEntry(key, &CRedisLLDPFields)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to collect LLDP info for interface %s", port[1])
	}
	capabilities, err := getCapabilities(rawData[CLLDPRemoteCapabilitiesSupported])
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode supported capabilities for remote interface")
	}
	enabledCapabilities, err := getCapabilities(rawData[CLLDPRemoteCapabilitiesEnabled])
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode enabled capabilities for remote interface")
	}

	f := &frame.Frame{
		InterfaceID:         fileVal,
		ChassisID:           rawData[CLLDPRemoteChassisId],
		SystemName:          rawData[CLLDPRemoteSystemName],
		SystemDescription:   rawData[CLLDPRemoteSystemDescription],
		Capabilities:        capabilities,
		EnabledCapabilities: enabledCapabilities,
		PortID:              rawData[CLLDPRemotePortId],
		PortDescription:     rawData[CLLDPRemotePortDescription],
		ManagementAddresses: strings.Split(rawData[CLLDPRemoteManagementAddresses], ","),
		TTL:                 0,
	}
	return f, nil
}

func getBitsList(num uint8) []int {
	bitsList := make([]int, 0)
	num = bits.Reverse8(num)
	for bit := 0; bit < 7; bit++ {
		if num&1 == 1 {
			bitsList = append(bitsList, bit)
		}
		num = num >> 1
	}
	return bitsList
}

func getCapabilities(caps string) ([]frame.Capability, error) {
	capabilities := make([]frame.Capability, 0)
	for _, i := range strings.Split(caps, " ") {
		if i == "00" || i == "" {
			continue
		}
		if parsed, err := strconv.ParseUint(i, 16, 8); err == nil {
			bitsList := getBitsList(uint8(parsed))
			for _, v := range bitsList {
				capabilities = append(capabilities, frame.CCapabilities[v])
			}
		} else {
			return nil, err
		}
	}
	return capabilities, nil
}
