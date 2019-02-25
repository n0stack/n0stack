package stringutil

import (
	"fmt"
	"hash/crc32"
	"math"
)

const SuffixLength = 5

func StringWithChecksumSuffix(name string, size int) string {
	cs := crc32.Checksum([]byte(name), crc32.IEEETable) % uint32(math.Pow(0x10, 4))

	if len(name)+SuffixLength > size {
		return fmt.Sprintf("%s-%02x", name[:size-SuffixLength], cs)
	}

	return fmt.Sprintf("%s-%02x", name, cs)
}
