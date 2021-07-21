package frame

import (
	"net"

	"github.com/pkg/errors"
)

func idBytesToMac(id []byte) (string, error) {
	idLen := len(id)
	if idLen != 6 {
		return "", errors.Errorf("expected to have 6 bytes in chassis ID, but got %d %v", idLen, id)
	}
	hwAddr := net.HardwareAddr(id)
	return hwAddr.String(), nil
}

func idBytesToNetworkAddress(id []byte) (string, error) {
	idLen := len(id)
	// According to RFC 2922 (page 11)
	// fist value octet is an identifier of an address type (v4 or v6)
	// the rest is a value itself
	validAddr := (idLen == 5 && id[0] == 1) || (idLen == 17 && id[0] == 2)
	if !validAddr {
		return "", errors.Errorf("expected to be IPv4 (1) or IPv6 address (2), but got %d with value %v", id[0], id[1:])
	}
	ip := net.IP(id[1:])
	return ip.String(), nil
}
