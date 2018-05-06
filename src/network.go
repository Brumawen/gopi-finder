package gopifinder

import (
	"net"
)

// GetLocalIPAddresses gets a list of valid IPv4 addresses for the local machine
func GetLocalIPAddresses() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	l := []string{}
	for _, i := range ifaces {
		adds, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range adds {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// Only select valid IPv4 addresses that are not loopbacks
			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
				l = append(l, ip.String())
			}
		}
	}
	return l, nil
}

// GetPotentialAddresses gets a list of IP addresses in the same subnet as the
// specified IP address that could, potentially, host a server.
// Note that this only supports Class 3 subnets for now.
func GetPotentialAddresses(ip string) ([]string, error) {
	a := net.ParseIP(ip).To4()
	l := []string{}
	if a != nil {
		for i := 2; i < 255; i++ {
			a[3] = byte(i)
			l = append(l, a.String())
		}
	}
	return l, nil
}
