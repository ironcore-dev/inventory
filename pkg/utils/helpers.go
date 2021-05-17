package utils

import (
	"math/bits"
	"os"

	"github.com/google/uuid"
)

const (
	CVersionFilePath = "/etc/sonic/sonic_version.yml"
	CMachineType     = "Machine"
	CSwitchType      = "Switch"
	CSonicNamespace  = "switch.onmetal.de"
)

func GetHostType(versionFile string) (string, error) {
	//todo: determining how to check host type without checking files
	if _, err := os.Stat(versionFile); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		} else {
			return CMachineType, nil
		}
	}
	return CSwitchType, nil
}

func GetUUID(namespace string, identifier string) string {
	namespaceUUID := uuid.NewMD5(uuid.UUID{}, []byte(namespace))
	newUUID := uuid.NewMD5(namespaceUUID, []byte(identifier))
	return newUUID.String()
}

func GetBitsList(num uint8, bitsList *[]int) {
	num = bits.Reverse8(num)
	for bit := 0; bit < 7; bit++ {
		if num&1 == 1 {
			*bitsList = append(*bitsList, bit)
		}
		num = num >> 1
	}
}
