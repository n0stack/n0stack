package nettools

import "net"

func NextIP(ip net.IP) net.IP {
	res := net.ParseIP(ip.String())
	for i := len(res) - 1; i >= 0; i-- {
		res[i]++

		if res[i] > 0 {
			break
		}
	}

	return res
}

func GetEndIP(ipn *net.IPNet) net.IP {
	// next network address
	ip := NextIP(ipn.IP)

	for {
		// before broadcast address
		if !ipn.Contains(NextIP(NextIP(ip))) {
			return ip
		}

		ip = NextIP(ip)
	}
}
