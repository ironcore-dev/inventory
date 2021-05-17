package lldp

import (
	"context"
	"math/bits"
	"path"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/file"
)

const (
	CLLDPEntryKeyMask = "LLDP_ENTRY*"
	CClassNetPath     = "/sys/class/net/"
	CIndexFile        = "ifindex"
)

var CRedisLLDPFields = []string{
	"lldp_rem_chassis_id",
	"lldp_rem_sys_name",
	"lldp_rem_sys_desc",
	"lldp_rem_sys_cap_supported",
	"lldp_rem_sys_cap_enabled",
	"lldp_rem_port_id",
	"lldp_rem_port_desc",
	"lldp_rem_man_addr",
}

type RedisSvc struct {
	client    *redis.Client
	ctx       context.Context
	indexPath string
}

func NewRedisSvc(basePath string) *RedisSvc {
	return &RedisSvc{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		ctx:       context.Background(),
		indexPath: path.Join(basePath, CClassNetPath),
	}
}

func (s *RedisSvc) GetFrames() ([]Frame, error) {
	frames := make([]Frame, 0)
	lldpKeys, err := s.GetKeysByPattern(CLLDPEntryKeyMask)
	if err != nil {
		return nil, err
	}
	for _, key := range lldpKeys {
		frame, err := s.processRedisPortData(key)
		if err != nil {
			return nil, err
		}
		frames = append(frames, *frame)
	}
	return frames, nil
}

func (s *RedisSvc) GetKeysByPattern(pattern string) ([]string, error) {
	val, err := s.client.Keys(s.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (s *RedisSvc) GetValuesFromHashEntry(key string, fields *[]string) (map[string]string, error) {
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

func (s *RedisSvc) processRedisPortData(key string) (*Frame, error) {
	port := strings.Split(key, ":")
	filePath := path.Join(s.indexPath, port[1], CIndexFile)
	fileVal, err := file.ToString(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get interface index value from %s", filePath)
	}
	rawData, err := s.GetValuesFromHashEntry(key, &CRedisLLDPFields)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to collect LLDP info for interface %s", port[1])
	}
	capabilities, err := getCapabilities(rawData["lldp_rem_sys_cap_supported"])
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode supported capabilities for remote interface")
	}
	enabledCapabilities, err := getCapabilities(rawData["lldp_rem_sys_cap_enabled"])
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode enabled capabilities for remote interface")
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

func getCapabilities(caps string) ([]Capability, error) {
	capabilities := make([]Capability, 0)
	for _, i := range strings.Split(caps, " ") {
		if i == "00" {
			continue
		}
		if parsed, err := strconv.ParseUint(i, 16, 8); err == nil {
			bitsList := getBitsList(uint8(parsed))
			for _, v := range bitsList {
				capabilities = append(capabilities, CCapabilities[v])
			}
		} else {
			return nil, err
		}
	}
	return capabilities, nil
}
