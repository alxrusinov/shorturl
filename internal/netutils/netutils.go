package netutils

import "net"

func CheckSubnet(trustedSubnet, ip string) (bool, error) {
	if trustedSubnet == "" {
		return false, nil
	}

	_, subnet, err := net.ParseCIDR(trustedSubnet)

	if err != nil {
		return false, err
	}

	ipNet := net.ParseIP(ip)

	trusted := subnet.Contains(ipNet)

	return trusted, nil
}
