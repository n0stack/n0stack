package kvm

import "strings"

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

// func (t tap) getMACAddr() *net.HardwareAddr {
// 	c := crc32.ChecksumIEEE(t.id.Bytes())

// 	return &net.HardwareAddr{0x52, 0x54, c[2:]}
// }
