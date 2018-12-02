package netutil

import (
	"fmt"
	"hash/crc32"
	"math"
)

const MaxLengthLinuxNetworkDeviceName = 15
const ChecksumLength = 5

// TrimNetdevName trim network device name because Linux network device can use 15 characters.
// コンフリクトを抑制するために末尾 4 bytes を乱数にする
func StructLinuxNetdevName(name string) string {
	cs := crc32.Checksum([]byte(name), crc32.IEEETable) % uint32(math.Pow(0x10, 4))

	if len(name)+ChecksumLength > MaxLengthLinuxNetworkDeviceName {
		return fmt.Sprintf("%s-%02x", name[:10], cs)
	}

	return fmt.Sprintf("%s-%02x", name, cs)
}
