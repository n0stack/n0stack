package netutil

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"net"
)

func GenerateHardwareAddress(id string) net.HardwareAddr {
	cs := crc32.Checksum([]byte(id), crc32.IEEETable)
	b, _ := hex.DecodeString(fmt.Sprintf("5254%08x", cs))
	return net.HardwareAddr(b)
}
