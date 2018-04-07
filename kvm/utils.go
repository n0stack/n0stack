package kvm

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"net"
	"strings"

	"github.com/satori/go.uuid"
)

func (k kvm) getInstanceName() string {
	i := strings.Split(k.id.String(), "-")
	return fmt.Sprintf("n0core-%s", i)
}

func (k kvm) getQMPPath() string {
	chardev := map[string]string{}
	chardevID := ""

	for i, a := range k.args {
		switch {
		case a == "-mon":
			ops := strings.Split(k.args[i+1], ",")
			for _, o := range ops {
				if strings.HasPrefix(o, "chardev=") {
					chardevID = strings.Split(o, "=")[1]
				}
			}

		case a == "-chardev":
			var (
				id string
				p  string
			)

			ops := strings.Split(k.args[i+1], ",")
			for _, o := range ops {
				switch {
				case strings.HasPrefix(o, "id="):
					id = strings.Split(o, "=")[1]
				case strings.HasPrefix(o, "path="):
					p = strings.Split(o, "=")[1]
				}
			}

			chardev[id] = p
		}
	}

	return chardev[chardevID]
}

func (k kvm) getMACAddr(id uuid.UUID) *net.HardwareAddr {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, crc32.ChecksumIEEE(id.Bytes()))
	s := fmt.Sprintf("52:54:%02x:%02x:%02x:%02x", b[0], b[1], b[2], b[3])
	h, _ := net.ParseMAC(s)

	return &h
}
