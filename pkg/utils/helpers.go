package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/google/uuid"
)

const (
	CVersionFilePath = "/etc/sonic/sonic_version.yml"
	CMachineType     = "Machine"
	CSwitchType      = "Switch"
)

func GetHostType() (string, error) {
	//todo: determining how to check host type without checking files
	if _, err := os.Stat(CVersionFilePath); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		} else {
			return CMachineType, nil
		}
	}
	return CSwitchType, nil
}

func GetUUID(text string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(text))
	rawUid := hex.EncodeToString(hasher.Sum(nil))
	if uid, err := uuid.Parse(rawUid); err == nil {
		return uid.String(), nil
	} else {
		return "", err
	}
}
