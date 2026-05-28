package xnet

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

// ParseAddr 解析地址
func ParseAddr(addr string) (listenAddr, exposeAddr string, err error) {
	var host, port string

	if addr != "" {
		host, port, err = net.SplitHostPort(addr)
		if err != nil {
			return
		}
	}

	if port == "" || port == "0" {
		p, err := AssignRandPort(host)
		if err != nil {
			return "", "", err
		}
		port = strconv.Itoa(p)
	}

	if len(host) > 0 && (host != "0.0.0.0" && host != "[::]" && host != "::") {
		listenAddr = net.JoinHostPort(host, port)
		exposeAddr = listenAddr
	} else {
		ip, err := InternalIP()
		if err != nil {
			return "", "", err
		}
		listenAddr = net.JoinHostPort("0.0.0.0", port)
		exposeAddr = net.JoinHostPort(ip, port)
	}

	return
}

// IP2Long IP地址转换为长整型
func IP2Long(ip string) uint32 {
	v := net.ParseIP(ip).To4()

	if len(v) == 0 {
		return 0
	}

	return binary.BigEndian.Uint32(v)
}

// Long2IP 长整型转换为字符串地址
func Long2IP(v uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, v)
	return ip.String()
}

// ExtractIP 提取主机地址
func ExtractIP(addr net.Addr) (ip string, err error) {
	ip, _, err = net.SplitHostPort(addr.String())
	return
}

// ExtractPort 提取主机端口
func ExtractPort(addr net.Addr) (int, error) {
	_, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(port)
}

// InternalIP 获取内网IP地址
func InternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var ip string
	for _, iface := range ifaces {
		if !isInternalIface(iface) {
			continue
		}

		next, preferred, err := ifaceInternalIP(iface)
		if err != nil {
			return "", err
		}
		if preferred {
			return next, nil
		}
		if ip == "" {
			ip = next
		}
	}

	if ip != "" {
		return ip, nil
	}
	return "", errors.New("not found ip address")
}

func isInternalIface(iface net.Interface) bool {
	return iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0
}

func ifaceInternalIP(iface net.Interface) (ip string, preferred bool, err error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return "", false, err
	}

	for _, addr := range addrs {
		ipnet, ok := addrIP(addr)
		if !ok {
			continue
		}

		ipv4 := ipnet.To4()
		if ipnet.IsLoopback() || ipv4 == nil || !ipv4.IsPrivate() {
			continue
		}
		if ipv4[0] == 192 && ipv4[1] == 168 {
			return ipv4.String(), true, nil
		}
		if ip == "" {
			ip = ipv4.String()
		}
	}

	return ip, false, err
}

func addrIP(addr net.Addr) (net.IP, bool) {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP, v.IP != nil
	case *net.IPAddr:
		return v.IP, v.IP != nil
	default:
		return nil, false
	}
}

// ExternalIP 获取外网IP地址
func ExternalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:54")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return ExtractIP(conn.LocalAddr())
}

// AssignRandPort 分配一个随机端口
func AssignRandPort(ip ...string) (int, error) {
	addr := ":0"
	if len(ip) > 0 {
		addr = ip[0] + addr
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port

	_ = listener.Close()

	return port, nil
}

// FulfillAddr 补全地址
func FulfillAddr(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	if host == "" {
		host = "0.0.0.0"
	}

	return net.JoinHostPort(host, port)
}
