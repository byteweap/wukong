package host

import (
	"fmt"
	"net"
	"strconv"
)

// ExtractHostPort from address
func ExtractHostPort(addr string) (host string, port uint64, err error) {
	var ports string
	host, ports, err = net.SplitHostPort(addr)
	if err != nil {
		return
	}
	port, err = strconv.ParseUint(ports, 10, 16) //nolint:mnd
	return
}

func isValidIP(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.IsGlobalUnicast() && !ip.IsInterfaceLocalMulticast()
}

// Port return a real port.
func Port(lis net.Listener) (int, bool) {
	if addr, ok := lis.Addr().(*net.TCPAddr); ok {
		return addr.Port, true
	}
	return 0, false
}

// Extract returns a private addr and port.
func Extract(hostPort string, lis net.Listener) (string, error) {
	addr, port, err := net.SplitHostPort(hostPort)
	if err != nil && lis == nil {
		return "", err
	}

	port, err = extractPort(port, lis)
	if err != nil {
		return "", err
	}
	if isSpecificAddr(addr) {
		return net.JoinHostPort(addr, port), nil
	}

	ip, err := extractPrivateIP()
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}
	return net.JoinHostPort(ip, port), nil
}

func extractPort(port string, lis net.Listener) (string, error) {
	if lis == nil {
		return port, nil
	}
	p, ok := Port(lis)
	if !ok {
		return "", fmt.Errorf("failed to extract port: %v", lis.Addr())
	}
	return strconv.Itoa(p), nil
}

func isSpecificAddr(addr string) bool {
	return len(addr) > 0 && addr != "0.0.0.0" && addr != "[::]" && addr != "::"
}

func extractPrivateIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var (
		minIndex = 0
		ips      = make([]net.IP, 0, 1)
	)
	for _, iface := range ifaces {
		ip := ifacePrivateIP(iface, minIndex, len(ips) != 0)
		if ip == nil {
			continue
		}
		minIndex = iface.Index
		ips = append(ips, ip)
	}
	if len(ips) != 0 {
		return ips[len(ips)-1].String(), nil
	}
	return "", nil
}

func ifacePrivateIP(iface net.Interface, minIndex int, found bool) net.IP {
	if iface.Flags&net.FlagUp == 0 || iface.Index >= minIndex && found {
		return nil
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil
	}
	var selected net.IP
	for _, rawAddr := range addrs {
		ip := rawAddrIP(rawAddr)
		if ip != nil && isValidIP(ip.String()) {
			selected = ip
			if ip.To4() != nil {
				return ip
			}
		}
	}
	return selected
}

func rawAddrIP(rawAddr net.Addr) net.IP {
	switch addr := rawAddr.(type) {
	case *net.IPAddr:
		return addr.IP
	case *net.IPNet:
		return addr.IP
	default:
		return nil
	}
}
