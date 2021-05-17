package utils

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

const (
	CLLDPEntryKeyMask = "LLDP_ENTRY*"
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

var ctx = context.Background()

func Client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetKeysByPattern(rdb *redis.Client, pattern string) ([]string, error) {
	val, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func GetValuesFromHashEntry(rdb *redis.Client, key string, fields *[]string) (map[string]string, error) {
	result := make(map[string]string)
	for _, f := range *fields {
		val, err := rdb.Do(ctx, "HGET", key, f).Result()
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
