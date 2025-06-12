package utils

import (
	"encoding/binary"
	"net"
)

func IpToUint32(ip net.IP) uint32 {
	ipv4 := ip.To4()
	return binary.BigEndian.Uint32(ipv4)
}
